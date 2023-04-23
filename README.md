# Filch

Filch shifted his career from being the caretaker of Hogwarts School of Witchcraft and Wizardry into the DevOps field, where he can take care of some manual tasks that you might want to apply to your cluster on a regular basis.
  

## Description

A for-fun Kubernetes Operator behaves similarly to the `CronJob` resource.
It was created for fun and to get more details about how the Operators are implemented

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster

Deploy the Operator
```sh
kubectl apply -f setup/crd-setup.yaml
```

Let's assume that you have a remote docker host whether it is deployed as pod inside your cluster or a remote instance and you want to delete any unused images which are older than 24 hours daily at 12:00AM. With the help of Mrs. Norris (Filch's cat) your desired task can be done as showing below

> Note: Docker Authentication is not considered in the example below, so be careful when you expose your docker api

1. Create docker-cleanup.yaml with the following content

```yaml
apiVersion: filch.caretaker.sh/v1
kind: MrsNorrisJob
metadata:
  labels:
    app.kubernetes.io/name: mrsnorrisjob
    app.kubernetes.io/instance: mrsnorrisjob-cleanup
    app.kubernetes.io/part-of: filch
    app.kubernetes.io/created-by: filch
  name: mrsnorrisjob-cleanup
spec:
  schedule: "0 * * * *"
  startingDeadlineSeconds: 60
  concurrencyPolicy: Forbid
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: docker-cleaner
            image: docker:cli
            args:
            - /bin/sh
            - -c
            - docker container prune --filter "until=24h" -f && docker image prune -a --filter "until=24h" -f
            env:
              - name: DOCKER_HOST
                value: tcp://my-remote-docker-engine:2375
          restartPolicy: OnFailure
```

2. Apply the file to your cluster
```sh
kubectl apply -f docker-cleanup.yaml
```

That's it!

## License

Distributed under the Apache License. See [`LICENSE`](LICENSE) for more information.

