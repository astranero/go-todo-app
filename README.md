# Hello Kubernetes

## Exercise 3.06: DBaaS vs DIY

 do-it-yourself (DIY) solution is appropriate when a specific database platform tailored to unique requirements is necessary. It’s ideal for situations where you have your own physical servers, which provide greater control over database design and configurations. However, this freedom comes with added complexity, as the database must be maintained and updated regularly.

As your services grow, more databases may be needed, which would require purchasing additional physical servers, increasing staffing for maintenance, managing electricity costs, and replacing old components every few years. You would also need to buy extra servers for failover scenarios, ensure data backups, and develop or set up effective monitoring and backup solutions. Despite these efforts, your data security is not guaranteed, especially if your infrastructure is localized in a single location. For example, a fire could destroy your servers, resulting in irreversible data loss. All of this requires significant time and financial investment, adding to organizational overhead, particularly when it comes to seeking permissions and approvals for database usage. Additionally, you would need to manage the security of your data independently.

In contrast, a Database as a Service (DBaaS) solution is much simpler in most cases. With DBaaS, you own a database without the need to design or manage the underlying software. The service provider handles the security, which is a significant advantage, as they typically take responsibility for any data loss. For instance, when a Microsoft update resulted in a costly issue for CrowdStrike, a company that handles Microsoft’s security. the impact on Microsoft itself was minimal because the responsibility was outsourced.

DBaaS solutions are usually ready to use, allowing you to quickly start with the specific requirements you need. Providers specializing in DBaaS invest substantial resources into developing robust security, backup, and other essential services for databases. Consequently, these solutions often offer a level of reliability and security that is difficult to achieve with a DIY approach.

## Exercise 5.08
![Exercise 5.08.png]
1. I used Flannel as the CNI (Container Network Interface) with k3d.
2. I have used Prometheus and Grafana for multiple purposes, including setting up a system at work with Thanos and other tools to monitor a multi-cluster solution, which includes monitoring GPUs.
3. I have configured Nginx as a proxy to forward requests to different API endpoints.
4. We use Calico to allocate IPs to pods.
5. I used Open vSwitch (OVS) to create virtual switches for virtual machines.
6. I deployed KubeVirt to create virtual machines within a Kubernetes cluster.
7. I utilized Ansible to configure virtual machines over a network.
8. I used MinIO for storing various backups, including Velero backups.
9. I have used Velero for backing up deployments.
10. We employed Longhorn as our storage solution.
11. I utilized containerd for container management, including namespace and process creation.
12. CoreDNS was used to handle Kubernetes pod name to IP address translation.
13. I used etcd to manage the state of the cluster.
14. Linkerd was employed in this course to perform canary releases, which depend on Envoy.
15. We used Harbor as a registry for storing container images in our cluster.
16. Traefik was set up to handle load balancing and ingress.
17. I used Multus to combine multiple CNIs (in our case, the pod network and OVS).
18. I have used Redis in an ML project for storing data in RAM.
19. I deployed Linstor as another storage solution for our nodes, which is particularly useful for diskless nodes.
20. I integrated Kaniko into our CI/CD pipeline to build and push images to Harbor.
21. I have used ArgoCD to automate repository pulls to the cluster in a GitOps manner.


## Exercise 5.05: Platform comparison - Rancer vs OpenShift

- Rancer offers an intuitive dashboard where users can easily install and manage applications for their cluster.
- OpenShift has a more complex interface with built-in development tools that facilitate testing and debugging.
- Rancer supports multiple orchestration platforms, including k8s and Docker swarm
- OpenShift focuses on K8s and integrates deeply with the RedHat.
- Rancher is know for its ease of installation and user-friendly interface, making it accessible for even those that have less k8s experience.
- OpenShift is more complex to setup and use, it requires a deeper understanding of the system and k8s.
- Rancer simplifies k8s operations through its UI, while OpenShifts rich feature set increases operational complexity.
- Rancer has a strong open-source community, with robust support & frequent updates.
- OpenShift has enterprise-level support, however with a cost.
- Rancer supports management of multiple clusters from differing cloud providers.
- Rancer is open-source and free, it is cost-effective choice when managing k8s clusters.
- OpenShift is a commercial product that has licensing costs, however it has enterprise lvel features and support.
- Rancer supports the management of multiple clusters across different cloud providers, it allows for multi-cloud environments.
- OpenShifts multi-cluster managements is tightlyy intergrated with RedHat's cloud and on-premise solutions.
- Rancer does not integrate as deeply with specific ecosystems as OpenShift, it focuses on broader compatibility.
- OpenShift integrates with RedHat system deeply, offering a lot of features for users within that environment.
- Rancer includes integration with Okta so that users can authenticate and be managed, it makes it secure and easy to use.
- OpenShifts security is included in its costly enterprise features.
