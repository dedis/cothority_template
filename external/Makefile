all: do-java

do-java:
	@for i in proto/*.proto ; \
	do \
		echo $$i ; \
		protoc -I=proto --java_out=java/src/main/java $$i ; \
	done
