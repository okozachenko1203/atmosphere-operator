package images

import (
	"fmt"
	"strings"

	dockerparser "github.com/novln/docker-parser"
)

var SKIP_LIST = []string{
	"image_repo_sync",
	"test",
}

var IMAGE_LIST = map[string]string{
	"alertmanager":                             "quay.io/prometheus/alertmanager:v0.24.0",
	"atmosphere":                               "quay.io/vexxhost/atmosphere:0.13.0", // x-release-please-version
	"barbican_api":                             "quay.io/vexxhost/barbican:wallaby",
	"barbican_db_sync":                         "quay.io/vexxhost/barbican:wallaby",
	"bootstrap":                                "quay.io/vexxhost/heat:zed",
	"ceph_config_helper":                       "quay.io/vexxhost/libvirtd:yoga-focal",
	"cert_manager_cainjector":                  "quay.io/jetstack/cert-manager-cainjector:v1.7.1",
	"cert_manager_cli":                         "quay.io/jetstack/cert-manager-ctl:v1.7.1",
	"cert_manager_controller":                  "quay.io/jetstack/cert-manager-controller:v1.7.1",
	"cert_manager_webhook":                     "quay.io/jetstack/cert-manager-webhook:v1.7.1",
	"cilium_node":                              "quay.io/cilium/cilium:v1.10.7@sha256:e23f55e80e1988db083397987a89967aa204ad6fc32da243b9160fbcea29b0ca",
	"cilium_operator":                          "quay.io/cilium/operator-generic:v1.10.7@sha256:d0b491d8d8cb45862ed7f0410f65e7c141832f0f95262643fa5ff1edfcddcafe",
	"cinder_api":                               "quay.io/vexxhost/cinder:zed",
	"cinder_backup_storage_init":               "quay.io/vexxhost/cinder:zed",
	"cinder_backup":                            "quay.io/vexxhost/cinder:zed",
	"cinder_db_sync":                           "quay.io/vexxhost/cinder:zed",
	"cinder_scheduler":                         "quay.io/vexxhost/cinder:zed",
	"cinder_storage_init":                      "quay.io/vexxhost/cinder:zed",
	"cinder_volume_usage_audit":                "quay.io/vexxhost/cinder:zed",
	"cinder_volume":                            "quay.io/vexxhost/cinder:zed",
	"cluster_api_controller":                   "registry.k8s.io/cluster-api/cluster-api-controller:v1.3.0",
	"cluster_api_kubeadm_bootstrap_controller": "registry.k8s.io/cluster-api/kubeadm-bootstrap-controller:v1.3.0",
	"cluster_api_kubeadm_control_plane_controller": "registry.k8s.io/cluster-api/kubeadm-control-plane-controller:v1.3.0",
	"cluster_api_openstack_controller":             "gcr.io/k8s-staging-capi-openstack/capi-openstack-controller:nightly_main_20221109",
	"csi_node_driver_registrar":                    "k8s.gcr.io/sig-storage/csi-node-driver-registrar:v2.4.0",
	"csi_rbd_attacher":                             "k8s.gcr.io/sig-storage/csi-attacher:v3.4.0",
	"csi_rbd_plugin":                               "quay.io/cephcsi/cephcsi:v3.5.1",
	"csi_rbd_provisioner":                          "k8s.gcr.io/sig-storage/csi-provisioner:v3.1.0",
	"csi_rbd_resizer":                              "k8s.gcr.io/sig-storage/csi-resizer:v1.3.0",
	"csi_rbd_snapshotter":                          "k8s.gcr.io/sig-storage/csi-snapshotter:v4.2.0",
	"db_drop":                                      "quay.io/vexxhost/heat:zed",
	"db_init":                                      "quay.io/vexxhost/heat:zed",
	"dep_check":                                    "quay.io/vexxhost/kubernetes-entrypoint:latest",
	"designate_api":                                "quay.io/vexxhost/designate:zed",
	"designate_central":                            "quay.io/vexxhost/designate:zed",
	"designate_db_sync":                            "quay.io/vexxhost/designate:zed",
	"designate_mdns":                               "quay.io/vexxhost/designate:zed",
	"designate_producer":                           "quay.io/vexxhost/designate:zed",
	"designate_sink":                               "quay.io/vexxhost/designate:zed",
	"designate_worker":                             "quay.io/vexxhost/designate:zed",
	"flux_helm_controller":                         "ghcr.io/fluxcd/helm-controller:v0.22.2",
	"flux_kustomize_controller":                    "ghcr.io/fluxcd/kustomize-controller:v0.27.0",
	"flux_notification_controller":                 "ghcr.io/fluxcd/notification-controller:v0.25.1",
	"flux_source_controller":                       "ghcr.io/fluxcd/source-controller:v0.26.1",
	"glance_api":                                   "quay.io/vexxhost/glance:zed",
	"glance_db_sync":                               "quay.io/vexxhost/glance:zed",
	"glance_metadefs_load":                         "quay.io/vexxhost/glance:zed",
	"glance_registry":                              "quay.io/vexxhost/glance:zed",
	"glance_storage_init":                          "quay.io/vexxhost/glance:zed",
	"grafana_sidecar":                              "quay.io/kiwigrid/k8s-sidecar:1.19.2",
	"grafana":                                      "docker.io/grafana/grafana:9.2.3",
	"haproxy":                                      "docker.io/library/haproxy:2.5",
	"heat_api":                                     "us-docker.pkg.dev/vexxhost-infra/openstack/heat:wallaby",
	"heat_cfn":                                     "us-docker.pkg.dev/vexxhost-infra/openstack/heat:wallaby",
	"heat_cloudwatch":                              "us-docker.pkg.dev/vexxhost-infra/openstack/heat:wallaby",
	"heat_db_sync":                                 "us-docker.pkg.dev/vexxhost-infra/openstack/heat:wallaby",
	"heat_engine_cleaner":                          "us-docker.pkg.dev/vexxhost-infra/openstack/heat:wallaby",
	"heat_engine":                                  "us-docker.pkg.dev/vexxhost-infra/openstack/heat:wallaby",
	"heat_purge_deleted":                           "us-docker.pkg.dev/vexxhost-infra/openstack/heat:wallaby",
	"horizon_db_sync":                              "us-docker.pkg.dev/vexxhost-infra/openstack/horizon:wallaby",
	"horizon":                                      "us-docker.pkg.dev/vexxhost-infra/openstack/horizon:wallaby",
	"ingress_nginx_controller":                     "k8s.gcr.io/ingress-nginx/controller:v1.1.1@sha256:0bc88eb15f9e7f84e8e56c14fa5735aaa488b840983f87bd79b1054190e660de",
	"ingress_nginx_default_backend":                "k8s.gcr.io/defaultbackend-amd64:1.5",
	"ingress_nginx_kube_webhook_certgen":           "k8s.gcr.io/ingress-nginx/kube-webhook-certgen:v1.1.1@sha256:64d8c73dca984af206adf9d6d7e46aa550362b1d7a01f3a0a91b20cc67868660",
	"keepalived":                                   "us-docker.pkg.dev/vexxhost-infra/openstack/keepalived:2.0.19",
	"keystone_api":                                 "quay.io/vexxhost/keystone:wallaby",
	"keystone_credential_cleanup":                  "quay.io/vexxhost/heat:zed",
	"keystone_credential_rotate":                   "quay.io/vexxhost/keystone:wallaby",
	"keystone_credential_setup":                    "quay.io/vexxhost/keystone:wallaby",
	"keystone_db_sync":                             "quay.io/vexxhost/keystone:wallaby",
	"keystone_domain_manage":                       "quay.io/vexxhost/heat:zed",
	"keystone_fernet_rotate":                       "quay.io/vexxhost/keystone:wallaby",
	"keystone_fernet_setup":                        "quay.io/vexxhost/keystone:wallaby",
	"ks_endpoints":                                 "quay.io/vexxhost/heat:zed",
	"ks_service":                                   "quay.io/vexxhost/heat:zed",
	"ks_user":                                      "quay.io/vexxhost/heat:zed",
	"kube_apiserver":                               "k8s.gcr.io/kube-apiserver:v1.22.17",
	"kube_controller_manager":                      "k8s.gcr.io/kube-controller-manager:v1.22.17",
	"kube_coredns":                                 "k8s.gcr.io/coredns/coredns:v1.8.4",
	"kube_etcd":                                    "k8s.gcr.io/etcd:3.5.0-0",
	"kube_proxy":                                   "k8s.gcr.io/kube-proxy:v1.22.17",
	"kube_scheduler":                               "k8s.gcr.io/kube-scheduler:v1.22.17",
	"kube_state_metrics":                           "registry.k8s.io/kube-state-metrics/kube-state-metrics:v2.6.0",
	"libvirt":                                      "quay.io/vexxhost/libvirtd:yoga-focal",
	"magnum_api":                                   "quay.io/vexxhost/magnum@sha256:46e7c910780864f4532ecc85574f159a36794f37aac6be65e4b48c46040ced17",
	"magnum_conductor":                             "quay.io/vexxhost/magnum@sha256:46e7c910780864f4532ecc85574f159a36794f37aac6be65e4b48c46040ced17",
	"magnum_db_sync":                               "quay.io/vexxhost/magnum@sha256:46e7c910780864f4532ecc85574f159a36794f37aac6be65e4b48c46040ced17",
	"magnum_registry":                              "docker.io/library/registry:2.7.1",
	"memcached":                                    "docker.io/library/memcached:1.6.17",
	"neutron_bagpipe_bgp":                          "us-docker.pkg.dev/vexxhost-infra/openstack/neutron:wallaby",
	"neutron_coredns":                              "docker.io/coredns/coredns:1.9.3",
	"neutron_db_sync":                              "us-docker.pkg.dev/vexxhost-infra/openstack/neutron:wallaby",
	"neutron_dhcp":                                 "us-docker.pkg.dev/vexxhost-infra/openstack/neutron:wallaby",
	"neutron_ironic_agent":                         "us-docker.pkg.dev/vexxhost-infra/openstack/neutron:wallaby",
	"neutron_l2gw":                                 "us-docker.pkg.dev/vexxhost-infra/openstack/neutron:wallaby",
	"neutron_l3":                                   "us-docker.pkg.dev/vexxhost-infra/openstack/neutron:wallaby",
	"neutron_linuxbridge_agent":                    "us-docker.pkg.dev/vexxhost-infra/openstack/neutron:wallaby",
	"neutron_metadata":                             "us-docker.pkg.dev/vexxhost-infra/openstack/neutron:wallaby",
	"neutron_netns_cleanup_cron":                   "us-docker.pkg.dev/vexxhost-infra/openstack/neutron:wallaby",
	"neutron_openvswitch_agent":                    "us-docker.pkg.dev/vexxhost-infra/openstack/neutron:wallaby",
	"neutron_server":                               "us-docker.pkg.dev/vexxhost-infra/openstack/neutron:wallaby",
	"neutron_sriov_agent_init":                     "us-docker.pkg.dev/vexxhost-infra/openstack/neutron:wallaby",
	"neutron_sriov_agent":                          "us-docker.pkg.dev/vexxhost-infra/openstack/neutron:wallaby",
	"node_feature_discovery":                       "k8s.gcr.io/nfd/node-feature-discovery:v0.11.2",
	"nova_api":                                     "quay.io/vexxhost/nova:wallaby",
	"nova_archive_deleted_rows":                    "quay.io/vexxhost/nova:wallaby",
	"nova_cell_setup_init":                         "quay.io/vexxhost/heat:zed",
	"nova_cell_setup":                              "quay.io/vexxhost/nova:wallaby",
	"nova_compute_ironic":                          "docker.io/kolla/ubuntu-source-nova-compute-ironic:wallaby",
	"nova_compute_ssh":                             "quay.io/vexxhost/nova-ssh:wallaby",
	"nova_compute":                                 "quay.io/vexxhost/nova:wallaby",
	"nova_conductor":                               "quay.io/vexxhost/nova:wallaby",
	"nova_consoleauth":                             "quay.io/vexxhost/nova:wallaby",
	"nova_db_sync":                                 "quay.io/vexxhost/nova:wallaby",
	"nova_novncproxy_assets":                       "quay.io/vexxhost/nova:wallaby",
	"nova_novncproxy":                              "quay.io/vexxhost/nova:wallaby",
	"nova_placement":                               "quay.io/vexxhost/nova:wallaby",
	"nova_scheduler":                               "quay.io/vexxhost/nova:wallaby",
	"nova_service_cleaner":                         "quay.io/vexxhost/cli:latest",
	"nova_spiceproxy_assets":                       "quay.io/vexxhost/nova:wallaby",
	"nova_spiceproxy":                              "quay.io/vexxhost/nova:wallaby",
	"octavia_api":                                  "quay.io/vexxhost/octavia:zed",
	"octavia_db_sync":                              "quay.io/vexxhost/octavia:zed",
	"octavia_health_manager_init":                  "quay.io/vexxhost/heat:zed",
	"octavia_health_manager":                       "quay.io/vexxhost/octavia:zed",
	"octavia_housekeeping":                         "quay.io/vexxhost/octavia:zed",
	"octavia_worker":                               "quay.io/vexxhost/octavia:zed",
	"openvswitch_db_server":                        "quay.io/vexxhost/openvswitch:2.17.3",
	"openvswitch_vswitchd":                         "quay.io/vexxhost/openvswitch:2.17.3",
	"percona_xtradb_cluster_haproxy":               "docker.io/percona/percona-xtradb-cluster-operator:1.10.0-haproxy",
	"percona_xtradb_cluster_operator":              "docker.io/percona/percona-xtradb-cluster-operator:1.10.0",
	"percona_xtradb_cluster":                       "docker.io/percona/percona-xtradb-cluster:5.7.39-31.61",
	"placement_db_sync":                            "quay.io/vexxhost/placement:wallaby",
	"placement":                                    "quay.io/vexxhost/placement:wallaby",
	"prometheus_config_reloader":                   "quay.io/prometheus-operator/prometheus-config-reloader:v0.60.1",
	"prometheus_ethtool_exporter":                  "quay.io/vexxhost/ethtool-exporter:5f05120a743a71adcbceb9f8ee1d43ecc7c4183a",
	"prometheus_ipmi_exporter":                     "us-docker.pkg.dev/vexxhost-infra/openstack/ipmi-exporter:1.4.0",
	"prometheus_memcached_exporter":                "quay.io/prometheus/memcached-exporter:v0.10.0",
	"prometheus_mysqld_exporter":                   "quay.io/prometheus/mysqld-exporter:v0.14.0",
	"prometheus_node_exporter":                     "quay.io/prometheus/node-exporter:v1.3.1",
	"prometheus_openstack_exporter":                "ghcr.io/openstack-exporter/openstack-exporter:1.6.0",
	"prometheus_operator_kube_webhook_certgen":     "k8s.gcr.io/ingress-nginx/kube-webhook-certgen:v1.3.0",
	"prometheus_operator":                          "quay.io/prometheus-operator/prometheus-operator:v0.60.1",
	"prometheus_pushgateway":                       "docker.io/prom/pushgateway:v1.4.2",
	"prometheus":                                   "quay.io/prometheus/prometheus:v2.39.1",
	"rabbit_init":                                  "docker.io/library/rabbitmq:3.10.2-management",
	"rabbitmq_cluster_operator":                    "docker.io/rabbitmqoperator/cluster-operator:1.13.1",
	"rabbitmq_credential_updater":                  "docker.io/rabbitmqoperator/default-user-credential-updater:1.0.2",
	"rabbitmq_server":                              "docker.io/library/rabbitmq:3.10.2-management",
	"rabbitmq_topology_operator":                   "docker.io/rabbitmqoperator/messaging-topology-operator:1.6.0",
	"senlin_api":                                   "us-docker.pkg.dev/vexxhost-infra/openstack/senlin:wallaby",
	"senlin_conductor":                             "us-docker.pkg.dev/vexxhost-infra/openstack/senlin:wallaby",
	"senlin_db_sync":                               "us-docker.pkg.dev/vexxhost-infra/openstack/senlin:wallaby",
	"senlin_engine_cleaner":                        "us-docker.pkg.dev/vexxhost-infra/openstack/senlin:wallaby",
	"senlin_engine":                                "us-docker.pkg.dev/vexxhost-infra/openstack/senlin:wallaby",
	"senlin_health_manager":                        "us-docker.pkg.dev/vexxhost-infra/openstack/senlin:wallaby",
	"skopeo":                                       "quay.io/skopeo/stable:latest",
	"tempest_run_tests":                            "us-docker.pkg.dev/vexxhost-infra/openstack/tempest:30.1.0-4",
}

func GetImageReference(imageName string) (*dockerparser.Reference, error) {
	if image, ok := IMAGE_LIST[imageName]; ok {
		ref, err := dockerparser.Parse(image)
		if err != nil {
			return nil, err
		}

		return ref, nil
	}

	return nil, fmt.Errorf("image not found: %s", imageName)
}

func OverrideRegistry(ref *dockerparser.Reference, registry string) (*dockerparser.Reference, error) {
	imageComponents := strings.Split(ref.ShortName(), "/")

	tag := ref.Tag()
	if strings.HasPrefix(tag, "sha256:") {
		tag = "@" + tag
	} else {
		tag = ":" + tag
	}

	image := fmt.Sprintf("%s/%s%s", registry, imageComponents[len(imageComponents)-1], tag)

	return dockerparser.Parse(image)
}
