# Keygen

OSS tool for easy SSH key generation

### Development Guide


#### Cluster Access
Login into cluster:

```bash
az aks get-credentials --resource-group keygen --name keygen
```

#### Bulding and publishing

To build and publish container:

(you need to login with quay.io to do that)

```
KEYGEN_VERSION=0.0.1 make build publish
```

Trigger redeploy:

```bash
./k8s/redeploy.sh
```


