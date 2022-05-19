package kubeclient

type KubeClient interface {
	kubectl
}

type kubeClient struct {
	kubectl
}

func NewKubeClient() (KubeClient, error) {
	k, err := newK()
	if err != nil {
		return nil, err
	}

	return &kubeClient{kubectl: k}, nil
}
