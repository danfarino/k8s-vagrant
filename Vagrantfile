script = <<-'END'
wget --no-verbose https://packages.cloud.google.com/apt/doc/apt-key.gpg
apt-key add apt-key.gpg
apt-add-repository "deb http://apt.kubernetes.io/ kubernetes-xenial main"
apt-get update
apt-get install kubeadm -y
sysctl net.bridge.bridge-nf-call-iptables=1

export MY_IP=$(ip addr | grep 'inet 192.168.' | awk '{ print $2 }' | sed -e 's/\(.*\)\/.*/\1/')
echo "KUBELET_EXTRA_ARGS='--node-ip=$MY_IP'" > /etc/default/kubelet
END

master_script = <<-'END'
kubeadm init --apiserver-advertise-address=192.168.11.101 --pod-network-cidr=10.244.0.0/16 | tee /root/kubeadm.log
# YOU MUST REMEMBER TO RUN THE NON-ROOT COMMANDS
export KUBECONFIG=/etc/kubernetes/admin.conf
#kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml
kubectl apply -f "https://cloud.weave.works/k8s/net?k8s-version=$(kubectl version | base64 | tr -d '\n')"
END

Vagrant.configure("2") do |config|
    config.vm.provider "virtualbox" do |v|
        v.cpus = 2
    end

    %w(master worker1 worker2).each do |name|
        config.vm.define name do |box|
            box.vm.box = "ubuntu/bionic64"  
            box.vm.hostname = name
            ip =
                case name
                when "master" then "192.168.11.101"
                when "worker1" then "192.168.11.102"
                when "worker2" then "192.168.11.103"
                end
            box.vm.network "private_network", ip: ip
            box.vm.provision "docker"
            box.vm.provision "shell", inline: script
            if name == "master" then
                box.vm.provision "shell", inline: master_script
            end
        end
    end
end
