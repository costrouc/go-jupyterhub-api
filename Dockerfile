FROM jupyterhub/jupyterhub

RUN pip install jupyterlab

RUN groupadd -g 1000 jupyter && \
    useradd jupyter -u 1000 -g 1000 -m -s /bin/bash

USER 1000:1000
ENV HOME=/home/jupyter
WORKDIR /home/jupyter

CMD ["jupyterhub"]