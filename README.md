# Test graceful termination in Kubernetes

Tests were run using this [deployment, service, and ingress](deployment.yaml).

Steps:

* Made sure two pods were running.
* Generated traffic
    ```
    $ ab -n 10000 -c 20 https://myhost.example.org/
    ```
* Deleted one of the pods with `kubectl delete pod <pod-name>`


### When keep-alive connections are allowed after SIGTERM

Active keep-alive connections will stay active until termination. New
connections are not established after the SIGTERM but requests are still
incomming on the active connections.

[keep_alive_enabled.log](keep_alive_enabled.log#L4082)


### When keep-alive connections are disabled after SIGTERM

keep-alive connections are closed and no new connections are established after
the SIGTERM.

[keep_alive_disabled_after_sigterm.log](keep_alive_disabled_after_sigterm.log#L852)
