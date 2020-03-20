## lokoctl cluster destroy

Destroy a cluster

### Synopsis

Destroy a cluster

```
lokoctl cluster destroy [flags]
```

### Options

```
      --confirm   Destroy cluster without asking for confirmation
  -h, --help      help for destroy
  -v, --verbose   Show output from Terraform
```

### Options inherited from parent commands

```
      --kubeconfig string     Path to kubeconfig file, taken from the asset dir if not given, and finally falls back to ~/.kube/config
      --lokocfg string        Path to lokocfg directory or file (default "./")
      --lokocfg-vars string   Path to lokocfg.vars file (default "./lokocfg.vars")
```

### SEE ALSO

* [lokoctl cluster](lokoctl_cluster.md)	 - Manage a cluster

