bin_path = os.getcwd()+'/bin'
if not os.environ.get('PATH').endswith(bin_path):
  os.environ.update({'PATH': os.environ.get('PATH')+':'+bin_path})

load('ext://kubebuilder', 'kubebuilder')
load('ext://min_k8s_version', 'min_k8s_version')

min_k8s_version('1.20')
kubebuilder('', 'opni.io', 'v1beta1', 'OpniCluster', 'rancher/opni')