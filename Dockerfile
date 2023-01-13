FROM alpine
ADD pod /pod
ADD filebeat.yml /filebeat.yml
ADD config /root/.kube/config
#Add filebeat /filebeat
ENTRYPOINT [ "/pod" ]