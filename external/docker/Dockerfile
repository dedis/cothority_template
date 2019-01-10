FROM dedis/conode_template:dev

# EXPOSE 7003 7005 7007 7009

COPY co1/private.toml co1/private.toml
COPY co2/private.toml co2/private.toml
COPY co3/private.toml co3/private.toml
COPY co4/private.toml co4/private.toml

COPY co1/public.toml co1/public.toml
COPY co2/public.toml co2/public.toml
COPY co3/public.toml co3/public.toml
COPY co4/public.toml co4/public.toml

RUN cp ./conode /usr/local/bin/
CMD ["./run_nodes.sh", "-v",  "2", "-n", "4" ]
