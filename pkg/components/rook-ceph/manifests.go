// Copyright 2020 The Lokomotive Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rookceph

// CephCluster resource definition was taken from https://github.com/rook/rook/blob/release-1.0/cluster/examples/kubernetes/ceph/cluster.yaml
const cephCluster = `
apiVersion: ceph.rook.io/v1
kind: CephCluster
metadata:
  name: rook-ceph
  namespace: {{ .Namespace }}
spec:
  cephVersion:
    image: ceph/ceph:v14.2.1-20190430
    allowUnsupported: true
  dataDirHostPath: /var/lib/rook
  mon:
    count: {{ .MonitorCount }}
    allowMultiplePerNode: false
  dashboard:
    enabled: true
  network:
    hostNetwork: false
  # RBD is required for block device. At the moment we only need object storage so this can be skipped.
  rbdMirroring:
    workers: 0
  placement:
    all:
      {{- if .NodeSelectors }}
      nodeAffinity:
        requiredDuringSchedulingIgnoredDuringExecution:
          nodeSelectorTerms:
            - matchExpressions:
              {{- range $item := .NodeSelectors }}
              - key: {{ $item.Key }}
                operator: {{ $item.Operator }}
                {{- if $item.Values }}
                values:
                  {{- range $val := $item.Values }}
                  - {{ $val }}
                  {{- end }}
                {{- end }}
              {{- end }}
      {{- end}}
      {{- if .TolerationsRaw }}
      tolerations: {{ .TolerationsRaw }}
      {{- end }}
  storage:
    useAllNodes: true
    useAllDevices: true
    config:
      {{- if .MetadataDevice }}
      metadataDevice: "{{ .MetadataDevice }}"
      {{- end }}
      storeType: bluestore
      osdsPerDevice: "1" # this value can be overridden at the node or device level
    # directories:
    # - path: /var/lib/rook
    #   # /dev/md/node-local-storage/rook
`
