# topz
A simple utility to expose top command as web interface.

Create a pod from the manifest

    kubectl apply -f mysql-pod-topz.yaml

Once the sidecar is running, check the pod IP address and visit `http://<pod_ip_addr>:1080/topz` to get a complete readout of the processes that are running in the mysql pod and their resource usage

    $ curl 172.17.0.10:1080/topz
    PID CPU % Memory %      Command
    1   0     0.00020143502 /pause
    7   0     0.12791124    /bin/bash /usr/local/bin/docker-entrypoint.sh mysqld
    78  0     0.22520435    /usr/bin/topz
    171 0     0.032833908   sleep 1

Credits to https://github.com/brendandburns/topz
