# imgcode
> PoC de embed de arquivos arbitrários em imagens e vídeos

Esse projeto de Abril de 2020 eu achei perdido no meu HD externo hoje (Janeiro de 2022).

Foi uma prova de conceito de embed de imagem e vídeo que pode ser recuperado depois.

Começou com um "Será que é possível usar aquele espaço ilimitado do Google Fotos pra guardar tranqueira arbitrária?" e a resposta é sim. 
Não cheguei a usar esse projeto de forma significativa e não faz diferença liberar pra geral porque eles acabaram com o storage ilimitado
do Google Fotos.

Fotos podiam ser recuperadas sem corromper se o formato for PNG e nenhum dos lados passar de 4000px.

Vídeos tinham limite de até 1080p e 10GB cada.

Em ambos os casos se passar do tamanho o Google Fotos quase sempre só recodifica pra um tamanho menor gerando perda de dados porque cada bit importa.

Vídeos foram mais complexos pois tinha que encodar num container que não faz compressão porque compressão de vídeo é lossy. 

Cada foto pode segurar coisa de 50MB.

Cada vídeo pode segurar coisa de alguns GB.

Tem na implementação uma forma de criptografar os dados usando uma senha. O overhead de foto é mínimo porque só tem os cabeçalhos do formatoo e o cabeçalho de metadados.

Vídeo é basicamente uma lista de imagens então primeiro são geradas as imagens e juntadas usando o ffmpeg.
