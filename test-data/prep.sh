#!/bin/bash

mkdir -p test-artefacts
curl https://repo1.maven.org/maven2/org/sonatype/maven-policy-demo/1.1.0/maven-policy-demo-1.1.0.jar --output test-artefacts/maven-policy-demo-1.1.0.jar
curl https://repo1.maven.org/maven2/org/sonatype/maven-policy-demo/1.1.0/maven-policy-demo-1.2.0.jar --output test-artefacts/maven-policy-demo-1.2.0.jar
curl https://repo1.maven.org/maven2/org/sonatype/maven-policy-demo/1.1.0/maven-policy-demo-1.3.0.jar --output test-artefacts/maven-policy-demo-1.3.0.jar
curl https://repo1.maven.org/maven2/com/amazonaws/aws-android-sdk-core/2.75.0/aws-android-sdk-core-2.75.0.aar --output test-artefacts/aws-android-sdk-core-2.75.0.aar
curl https://repo1.maven.org/maven2/org/jsoup/jsoup/1.13.1/jsoup-1.13.1.jar --output test-artefacts/jsoup-1.13.1.jar
curl https://repo1.maven.org/maven2/ant/ant/1.6.5/ant-1.6.5.jar --output test-artefacts/ant-1.6.5.jar
curl https://repo1.maven.org/maven2/org/springframework/spring-context/6.2.3/spring-context-6.2.3.jar --output test-artefacts/spring-context-6.2.3.jar
