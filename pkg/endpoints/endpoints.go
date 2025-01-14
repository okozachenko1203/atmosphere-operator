package endpoints

import (
	"fmt"

	"golang.org/x/exp/slices"
	"helm.sh/helm/v3/pkg/chart"
)

func endpointAuthUsers(endpoint interface{}) []string {
	skipList := []string{"test"}
	keys := []string{}

	for k := range endpoint.(map[string]interface{})["auth"].(map[string]interface{}) {
		if slices.Contains(skipList, k) {
			continue
		}

		keys = append(keys, k)
	}
	return keys
}

func endpointAuth(endpoint interface{}, passwords map[string]string) (map[string]interface{}, error) {
	auth := map[string]interface{}{}

	for _, user := range endpointAuthUsers(endpoint) {
		password, ok := passwords[user]
		if !ok {
			return nil, fmt.Errorf("password for user %s missing", user)
		}

		auth[user] = map[string]interface{}{
			"password": password,
		}
	}

	return auth, nil
}

func ForChart(chart *chart.Chart, config *EndpointConfig) (map[string]interface{}, error) {
	endpoints := make(map[string]interface{})
	chartEndpoints := chart.Values["endpoints"].(map[string]interface{})

	for endpointName, endpointValues := range chartEndpoints {
		switch endpointName {
		case "cluster_domain_suffix":
		case "ingress":
		case "kube_dns":
		case "ldap":
		case "local_image_registry":
		case "oci_image_registry":
		case "fluentd":
			continue
		case "oslo_cache":
			// TODO
			endpoints[endpointName] = map[string]interface{}{
				"auth": map[string]interface{}{
					"memcache_secret_key": config.MemcacheSecretKey,
				},
			}
		case "oslo_db":
			auth, err := endpointAuth(endpointValues, map[string]string{
				"admin":    config.DatabaseRootPassword,
				"keystone": config.KeystoneDatabasePassword,
			})
			if err != nil {
				return nil, err
			}

			endpoints[endpointName] = map[string]interface{}{
				"auth":      auth,
				"namespace": config.DatabaseNamespace,
				"hosts": map[string]interface{}{
					"default": config.DatabaseServiceName,
				},
			}
		case "oslo_messaging":
			auth, err := endpointAuth(endpointValues, map[string]string{
				"admin":    config.RabbitmqAdminPassword,
				"keystone": config.KeystoneRabbitmqPassword,
			})
			if err != nil {
				return nil, err
			}

			// NOTE(mnaser): RabbitMQ operator generates random usernames for the
			//               admin user.
			auth["admin"].(map[string]interface{})["username"] = config.RabbitmqAdminUsername

			endpoints[endpointName] = map[string]interface{}{
				"auth":        auth,
				"statefulset": nil,
				"namespace":   config.RabbitmqNamespace,
				"hosts": map[string]interface{}{
					"default": config.RabbitmqServiceName,
				},
			}
		case "identity":
			auth, err := endpointAuth(endpointValues, map[string]string{
				"admin": config.KeystoneAdminPassword,
			})
			if err != nil {
				return nil, err
			}

			for _, user := range endpointAuthUsers(endpointValues) {
				auth[user].(map[string]interface{})["region_name"] = config.RegionName
				auth[user].(map[string]interface{})["username"] = fmt.Sprintf("%s-%s", user, config.RegionName)
			}

			endpoints[endpointName] = map[string]interface{}{
				"auth": auth,
				"hosts": map[string]interface{}{
					"default": "keystone-api",
				},
				"scheme": map[string]interface{}{
					"public": "https",
				},
				"host_fqdn_override": map[string]interface{}{
					"public": map[string]interface{}{
						"host": config.KeystoneHost,
					},
				},
				"port": map[string]interface{}{
					"api": map[string]interface{}{
						"default": 5000,
						"public":  443,
					},
				},
			}
		default:
			return nil, fmt.Errorf("endpoint %s not supported", endpointName)
		}
	}

	return endpoints, nil
}
