# GoCam

- Tem endpoint em `/record` que inicializa o ffmpeg para gravar a camera para um ficheiro .mp4 (o nome do ficheiro é a data e hora do inicio da gravaçao)
  - com o comando start começa a gravar
  - com o comando stop para de gravar
	- ao parar de gravar copia o ficheiro mp4 para o `/mnt/pve/external/camera_output/`
  - caso não haja o comando stop, a gravação terminada passado 1 hora
  - a ideia é colocar também um timer, ou seja, fazer o comando start e dizer `start 2h` e ele fica 2h a gravar
