#!/bin/bash

kubernetes_version=${kubernetes_version:=v1.24.6}
crictl_release_version=${crictl_release_version:=v1.24.6}
kubernetes_release_version=${kubernetes_release_version:=v0.10.0}
containerd_version=${containerd_version:=1.6.8}
runc_version=${runc_version:=1.1.4}
ARCH=${ARCH:="amd64"}

echo "Node setup starting"

# Required on Scaleway CentOS machines
sudo bash -c 'cat > /etc/modules-load.d/containerd.conf' << EOF
overlay
br_netfilter
EOF

sudo modprobe overlay
sudo modprobe br_netfilter

# Setup required sysctl params, these persist across reboots.
sudo bash -c 'cat > /etc/sysctl.d/99-kubernetes-cri.conf' << EOF
net.bridge.bridge-nf-call-iptables  = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.ip_forward                 = 1
EOF

sudo sysctl --system

# disable predictable interface names
# sudo sed -i.bak 's/GRUB_CMDLINE_LINUX_DEFAULT=.*/GRUB_CMDLINE_LINUX_DEFAULT="net.ifnames=0 bios.devname=0 quiet"/' /etc/default/grub
# sudo update-grub


# Install containerd
## Set up the repository
echo "Installing required packages"
sudo yum install yum-utils device-mapper-persistent-data lvm2 gnupg2 conntrack libseccomp -y

echo "Adding docker repository"
sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo

echo "Updating"
# sudo yum update -y --exclude=kubelet,kubeadm,containerd,docker,kubectl,kernel

echo "installing containerd"

# sudo yum install containerd.io-${containerd_version} -y
curl https://github.com/containerd/containerd/releases/download/v${containerd_version}/cri-containerd-cni-${containerd_version}-linux-amd64.tar.gz -LJO
sudo tar --no-overwrite-dir -C / -xzf cri-containerd-cni-${containerd_version}-linux-amd64.tar.gz

# clean up the downloaded file
rm -rf cri-containerd-cni-${containerd_version}-linux-amd64.tar.gz

# Delete containerd cni configuration, as we depend on either calico or weavenet CNI
sudo rm -f /etc/cni/net.d/10-containerd-net.conflist
# sudo rm -rf /opt/cni

echo "Configuring containerd"
sudo mkdir -p /etc/containerd
sudo touch /etc/containerd/config.toml

sudo bash -c 'cat > /etc/containerd/config.toml << EOF
disabled_plugins = []
imports = []
oom_score = 0
plugin_dir = ""
required_plugins = []
root = "/var/lib/containerd"
state = "/run/containerd"
temp = ""
version = 2

[cgroup]
  path = ""

[debug]
  address = ""
  format = ""
  gid = 0
  level = ""
  uid = 0

[grpc]
  address = "/var/run/containerd/containerd.sock"
  gid = 0
  max_recv_message_size = 16777216
  max_send_message_size = 16777216
  tcp_address = ""
  tcp_tls_ca = ""
  tcp_tls_cert = ""
  tcp_tls_key = ""
  uid = 0

[metrics]
  address = ""
  grpc_histogram = false

[plugins]

  [plugins."io.containerd.gc.v1.scheduler"]
    deletion_threshold = 0
    mutation_threshold = 100
    pause_threshold = 0.02
    schedule_delay = "0s"
    startup_delay = "100ms"

  [plugins."io.containerd.grpc.v1.cri"]
    device_ownership_from_security_context = false
    disable_apparmor = false
    disable_cgroup = false
    disable_hugetlb_controller = true
    disable_proc_mount = false
    disable_tcp_service = true
    enable_selinux = false
    enable_tls_streaming = false
    enable_unprivileged_icmp = false
    enable_unprivileged_ports = false
    ignore_image_defined_volumes = false
    max_concurrent_downloads = 3
    max_container_log_line_size = 16384
    netns_mounts_under_state_dir = false
    restrict_oom_score_adj = false
    sandbox_image = "k8s.gcr.io/pause:3.6"
    selinux_category_range = 1024
    stats_collect_period = 10
    stream_idle_timeout = "4h0m0s"
    stream_server_address = "127.0.0.1"
    stream_server_port = "0"
    systemd_cgroup = false
    tolerate_missing_hugetlb_controller = true
    unset_seccomp_profile = ""

    [plugins."io.containerd.grpc.v1.cri".cni]
      bin_dir = "/opt/cni/bin"
      conf_dir = "/etc/cni/net.d"
      conf_template = ""
      ip_pref = ""
      max_conf_num = 1

    [plugins."io.containerd.grpc.v1.cri".containerd]
      default_runtime_name = "runc"
      disable_snapshot_annotations = true
      discard_unpacked_layers = false
      ignore_rdt_not_enabled_errors = false
      no_pivot = false
      snapshotter = "overlayfs"

      [plugins."io.containerd.grpc.v1.cri".containerd.default_runtime]
        base_runtime_spec = ""
        cni_conf_dir = ""
        cni_max_conf_num = 0
        container_annotations = []
        pod_annotations = []
        privileged_without_host_devices = false
        runtime_engine = ""
        runtime_path = ""
        runtime_root = ""
        runtime_type = ""

        [plugins."io.containerd.grpc.v1.cri".containerd.default_runtime.options]

      [plugins."io.containerd.grpc.v1.cri".containerd.runtimes]

        [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc]
          base_runtime_spec = ""
          cni_conf_dir = ""
          cni_max_conf_num = 0
          container_annotations = []
          pod_annotations = []
          privileged_without_host_devices = false
          runtime_engine = ""
          runtime_path = ""
          runtime_root = ""
          runtime_type = "io.containerd.runc.v2"

          [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc.options]
            BinaryName = ""
            CriuImagePath = ""
            CriuPath = ""
            CriuWorkPath = ""
            IoGid = 0
            IoUid = 0
            NoNewKeyring = false
            NoPivotRoot = false
            Root = ""
            ShimCgroup = ""
            SystemdCgroup = false

      [plugins."io.containerd.grpc.v1.cri".containerd.untrusted_workload_runtime]
        base_runtime_spec = ""
        cni_conf_dir = ""
        cni_max_conf_num = 0
        container_annotations = []
        pod_annotations = []
        privileged_without_host_devices = false
        runtime_engine = ""
        runtime_path = ""
        runtime_root = ""
        runtime_type = ""

        [plugins."io.containerd.grpc.v1.cri".containerd.untrusted_workload_runtime.options]

    [plugins."io.containerd.grpc.v1.cri".image_decryption]
      key_model = "node"

    [plugins."io.containerd.grpc.v1.cri".registry]
      config_path = ""

      [plugins."io.containerd.grpc.v1.cri".registry.auths]

      [plugins."io.containerd.grpc.v1.cri".registry.configs]

      [plugins."io.containerd.grpc.v1.cri".registry.headers]

      [plugins."io.containerd.grpc.v1.cri".registry.mirrors]

    [plugins."io.containerd.grpc.v1.cri".x509_key_pair_streaming]
      tls_cert_file = ""
      tls_key_file = ""

  [plugins."io.containerd.internal.v1.opt"]
    path = "/opt/containerd"

  [plugins."io.containerd.internal.v1.restart"]
    interval = "10s"

  [plugins."io.containerd.internal.v1.tracing"]
    sampling_ratio = 1.0
    service_name = "containerd"

  [plugins."io.containerd.metadata.v1.bolt"]
    content_sharing_policy = "shared"

  [plugins."io.containerd.monitor.v1.cgroups"]
    no_prometheus = false

  [plugins."io.containerd.runtime.v1.linux"]
    no_shim = false
    runtime = "runc"
    runtime_root = ""
    shim = "containerd-shim"
    shim_debug = false

  [plugins."io.containerd.runtime.v2.task"]
    platforms = ["linux/amd64"]
    sched_core = false

  [plugins."io.containerd.service.v1.diff-service"]
    default = ["walking"]

  [plugins."io.containerd.service.v1.tasks-service"]
    rdt_config_file = ""

  [plugins."io.containerd.snapshotter.v1.aufs"]
    root_path = ""

  [plugins."io.containerd.snapshotter.v1.btrfs"]
    root_path = ""

  [plugins."io.containerd.snapshotter.v1.devmapper"]
    async_remove = false
    base_image_size = ""
    discard_blocks = false
    fs_options = ""
    fs_type = ""
    pool_name = ""
    root_path = ""

  [plugins."io.containerd.snapshotter.v1.native"]
    root_path = ""

  [plugins."io.containerd.snapshotter.v1.overlayfs"]
    root_path = ""
    upperdir_label = false

  [plugins."io.containerd.snapshotter.v1.zfs"]
    root_path = ""

  [plugins."io.containerd.tracing.processor.v1.otlp"]
    endpoint = ""
    insecure = false
    protocol = ""

[proxy_plugins]

[stream_processors]

  [stream_processors."io.containerd.ocicrypt.decoder.v1.tar"]
    accepts = ["application/vnd.oci.image.layer.v1.tar+encrypted"]
    args = ["--decryption-keys-path", "/etc/containerd/ocicrypt/keys"]
    env = ["OCICRYPT_KEYPROVIDER_CONFIG=/etc/containerd/ocicrypt/ocicrypt_keyprovider.conf"]
    path = "ctd-decoder"
    returns = "application/vnd.oci.image.layer.v1.tar"

  [stream_processors."io.containerd.ocicrypt.decoder.v1.tar.gzip"]
    accepts = ["application/vnd.oci.image.layer.v1.tar+gzip+encrypted"]
    args = ["--decryption-keys-path", "/etc/containerd/ocicrypt/keys"]
    env = ["OCICRYPT_KEYPROVIDER_CONFIG=/etc/containerd/ocicrypt/ocicrypt_keyprovider.conf"]
    path = "ctd-decoder"
    returns = "application/vnd.oci.image.layer.v1.tar+gzip"

[timeouts]
  "io.containerd.timeout.bolt.open" = "0s"
  "io.containerd.timeout.shim.cleanup" = "5s"
  "io.containerd.timeout.shim.load" = "5s"
  "io.containerd.timeout.shim.shutdown" = "3s"
  "io.containerd.timeout.task.state" = "2s"

[ttrpc]
  address = ""
  gid = 0
  uid = 0
EOF
'

if [ -z "${runc_version}" ]
then
  echo "not installing custom runc version"
  echo "using runc versioned with containerd"
else
  RUNC_VERSION="v${runc_version}"
  echo "installing custom runc version: ${RUNC_VERSION}"
  curl https://github.com/opencontainers/runc/releases/download/${RUNC_VERSION}/runc.${ARCH} -LJO
  sudo install -m 755 runc.${ARCH} /usr/local/sbin/runc
fi
/usr/local/sbin/runc --version

echo "Enabling containerd"
sudo systemctl enable containerd
echo "Restarting containerd"
sudo systemctl restart containerd

echo "Create download DIR"
sudo mkdir -p /usr/local/bin

# Add /usr/local/bin to path
sudo bash -c "cat > /etc/profile.d/usr-local-bin.sh << EOF
#!/bin/bash

PATH=/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin:/root/bin
export PATH
EOF
"

# Install crictl
echo "Set crictl release version"
export CRICTL_VERSION="${crictl_release_version}"

echo "Download crictl"
sudo curl https://github.com/kubernetes-sigs/cri-tools/releases/download/${CRICTL_VERSION}/crictl-${CRICTL_VERSION}-linux-amd64.tar.gz -LJO
sudo tar -C /usr/local/bin -xzf crictl-${CRICTL_VERSION}-linux-amd64.tar.gz

# clean up the downloaded file
rm -rf crictl-${CRICTL_VERSION}-linux-amd64.tar.gz

# import google cloud packages GPG keys
echo "Import GPG keys"
sudo rpm --import https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg

echo "Setting SELinux to permissive mode (effectively disabling it)"
# Setting SELinux in permissive mode by running setenforce 0 and sed ... effectively disables it.
# This is required to allow containers to access the host filesystem, which is needed by pod networks for example.
# You have to do this until SELinux support is improved in the kubelet.
# setenforce 0
sudo sed -i 's/^SELINUX=enforcing$/SELINUX=permissive/' /etc/selinux/config

echo "Installing Kubelet, Kubeadm and Kubectl"

rm /usr/local/bin/{kubeadm,kubelet,kubectl} || true
cd /usr/local/bin
sudo curl -L --remote-name-all https://storage.googleapis.com/kubernetes-release/release/${kubernetes_version}/bin/linux/amd64/{kubeadm,kubelet,kubectl}
sudo chmod +x {kubeadm,kubelet,kubectl}

# RELEASE_VERSION="${kubernetes_release_version}"
curl -sSL "https://raw.githubusercontent.com/kubernetes/release/${kubernetes_release_version}/cmd/kubepkg/templates/latest/rpm/kubelet/kubelet.service" | sed "s:/usr/bin:/usr/local/bin:g" | sudo tee /etc/systemd/system/kubelet.service

sudo mkdir -p /etc/systemd/system/kubelet.service.d
curl -sSL "https://raw.githubusercontent.com/kubernetes/release/${kubernetes_release_version}/cmd/kubepkg/templates/latest/rpm/kubeadm/10-kubeadm.conf" | sed "s:/usr/bin:/usr/local/bin:g" | sudo tee /etc/systemd/system/kubelet.service.d/10-kubeadm.conf

mkdir -p /var/lib/kubelet
echo 'KUBELET_KUBEADM_ARGS="--container-runtime=remote --container-runtime-endpoint=/run/containerd/containerd.sock"' > /var/lib/kubelet/kubeadm-flags.env

echo "Enabling kubelet"
sudo systemctl enable --now -f kubelet
sudo systemctl restart kubelet

echo "Enabling wireguard"
# Install wireguard (https://www.wireguard.com/install/) for calico network encryption.
# No further configuration is necessary, as this is handled by calico.
sudo yum install epel-release elrepo-release -y
sudo yum install yum-plugin-elrepo -y
sudo yum install kmod-wireguard wireguard-tools -y

# Change journald memory settings (AME-725)
echo "Change journald memory settings (AME-725)"
sudo sed -i 's/#RuntimeMaxUse=/RuntimeMaxUse=10M/' /etc/systemd/journald.conf

echo "Node setup completed"

exit 0;
