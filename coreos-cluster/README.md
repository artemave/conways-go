CoreOs based game servers cluster
=========

## To run in Vagrant ##

This setup is largely copied from [here](http://blog.dixo.net/2015/02/load-balancing-with-coreos/). Nginx configuration is slightly different so that it reroutes on 404 as well (http://serverfault.com/questions/481454/load-balancer-to-handle-server-errors-silently).

    % vargant up                                # starts coreos cluster (3 servers)
    % export FLEETCTL_TUNNEL=127.0.0.1:2222     # so that fleetctl knows to talk to vagrant
    % ssh-add ~/.vagrant.d/insecure_private_key # might need that too for fleetctl <-> vagrant integration
    % fleetctl load instances/*                 # load units onto machines
    % fleetctl start instances/*                # start them

The above relies on [artemave/conways-go](https://registry.hub.docker.com/u/artemave/conways-go/) and [artemave/confd](https://registry.hub.docker.com/u/artemave/confd/) being publicly available in docker hub.
