---
name: Support
about: If you have questions about reaper-operator
title: ''
labels: needs-triage, question
assignees: ''

---

<!--
Thanks for filing an issue! Before hitting the button, please answer these questions.

Fill in as much of the template below as you can.

Note that this repository is about reaper-operator itself. If you have questions specifically about Reaper, please visit https://github.com/thelastpickle/cassandra-reaper.

We will try our best to answer the question, but we also have a mailing list (k8ssandra-users@googlegroups.com.) for any other questions.
-->

**Type of question**
<!-- Uncomment one or more of the following lines depending on what you are asking about: -->

- [ ] Best practices
- [ ] How to perform a particular operation
- [ ] Cassandra-related question
- [ ] Monitoring-related question
- [ ] Repair-related question
- [ ] Backup/restore-related question
- [ ] Open question

**What did you do?**

**Did you expect to see some different?**

**Environment (please complete the following information):**

* reaper-operator version:
<!-- Insert the image tag or Git SHA here. -->

<!--
    You can try a jsonpath query with kubectl like this to get the version:

        kubectl get deployment <reaper-operator-deployment> \
            -o jsonpath='{.spec.template.spec.containe[0].image}'
-->

<!--
Please provide the follow info if you deployed reaper-operator via the
k8ssandra Helm chart(s). 
-->
* Helm charts version info 
<!-- list installed charts and their versions from all namespaces -->
<!-- Replace the command with its output -->
`$ helm ls -A` 

* Helm charts user-supplied values
<!-- For each k8ssandra chart involved list user-supplied values -->
<!-- Replace the commands with its output -->
`$ helm get values RELEASE_NAME` 

* Kubernetes version information:
<!-- Replace the command with its output -->
`kubectl version`

* Kubernetes cluster kind:
<!-- Insert how you created your cluster: kind, kops, bootkube, etc. -->

* Manifests:

<!-- Please provide any manifests relevant to the issue -->

* Operator logs:

<!-- Please provide any reaper-operator logs relevant to the issue -->
<!-- 
  You can try a command like the following to get the logs if the operator was
  installed with Helm.
 -->
`kubectl -n <namespace> logs -l name=<releaseName>-reaper-operator-k8ssandra`

**Additional context**
<!-- Add any other context about the problem here. -->
