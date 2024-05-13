University: [ITMO University](https://itmo.ru/ru/)  
Faculty: [FICT](https://fict.itmo.ru)  
Course: [Application containerization and orchestration](https://github.com/itmo-ict-faculty/application-containerization-and-orchestration)  
Year: 2023/2024  
Group: K4111c  
Author: Filippov Artem Alekseevich  
Practice: practice2  
Date of create: 01.05.2024  
Date of finished: 03.05.2024  

Цель: изучить принципы работы с базами данных в контексте микросервисных приложений

Ход работы:

1. Выбрана и развернута база данных PostgreSQL в кластере minikube.

    ```yaml
    ---
    apiVersion: v1
    kind: PersistentVolume
    metadata:
      name: postgres-pv
    spec:
      accessModes:
        - ReadWriteOnce
      capacity:
        storage: 2Gi
      hostPath:
        path: /data/postgres
      storageClassName: standard
    ---
    apiVersion: v1
    kind: PersistentVolumeClaim
    metadata:
      name: postgres-pvc
    spec:
      accessModes:
        - ReadWriteOnce
      resources:
        requests:
        storage: 2Gi
      storageClassName: standard
      volumeName: postgres-pv
    ---
    apiVersion: apps/v1
    kind: StatefulSet
    metadata:
      name: postgres
    labels:
      service: postgres
    spec:
      replicas: {{ .Values.db.replicas }}
      selector:
        matchLabels:
          service: postgres
      serviceName: postgres
      template:
        metadata:
        labels:
          service: postgres
        spec:
          containers:
            - name: postgres
            image: postgres:latest
            imagePullPolicy: Always
            ports:
              - containerPort: {{ .Values.db.port }}
                name: postgres
                protocol: TCP
            env:
              - name: POSTGRES_DB
                value: {{ .Values.db.name }}
              - name: POSTGRES_USER
                value: {{ .Values.db.user }}
              - name: POSTGRES_PASSWORD
                valueFrom:
                  secretKeyRef:
                    name: postgres-secret
                    key: postgres-password
            volumeMounts:
              - mountPath: /var/lib/postgresql/data
                name: postgres-pvc
          volumes:
            - name: postgres-pvc
              persistentVolumeClaim:
              claimName: postgres-pvc
    ---
    apiVersion: v1
    kind: Service
    metadata:
      name: postgres
    labels:
      service: postgres
    spec:
      selector:
        service: postgres
      clusterIP: None
      ports:
        - name: http
          port: {{ .Values.db.port }}
          targetPort: {{ .Values.db.port }}
          protocol: TCP
    ```

Вывод: в ходе выполнения практической работы была выбрана база данных PostgreSQL, после чего был написан манифест для ее развертывания в кластере minikube с последующим развертыванием.
