# golang websocket-server and webassembly-client  service-framework

go 언어로 만든 websocket-server, webassembly-client 를 지원하는 서비스 프레임워크 입니다. 

goguelike, gowasm2dgame, gowasm3dgame, gofieldwar(비공개) 등을 만들며 사용한 코드들의 
공통 부분들을 모아서 다른 프로젝트에도 사용하기 편하게 정리한 것입니다. 

기본 코드는 stage 기반의 채팅기능만이 존재하며 각 프로젝트의 필요에 따라 확장해 쓰시면 됩니다. 

html5 canvas 를 사용한 2D game 예제를 보려면 https://github.com/kasworld/gowasm2dgame

webgl 을 사용한 3D game 예제를 보려면 https://github.com/kasworld/gowasm3dgame

좀 복잡한 게임(수만 라인 규모) 예제를 보려면 https://github.com/kasworld/goguelike

를 참고하시면 됩니다. 


# 사전 준비 사항 ( goguelike 의 INSTALL.md 참고)

준비물 : linux(debian,ubuntu,mint) , chrome web brower , golang 

버전 string 생성시 사용 : sha256sum, awk

goimports : 소스 코드 정리, import 해결

    go get golang.org/x/tools/cmd/goimports

프로토콜 생성기 : https://github.com/kasworld/genprotocol

    go get github.com/kasworld/genprotocol

Enum 생성기 : https://github.com/kasworld/genenum

    go get github.com/kasworld/genenum

Log 패키지 및 커스텀 로그레벨 로거 생성기 : https://github.com/kasworld/log

    go get github.com/kasworld/log
    install.sh 실행해서 genlog 생성 

# 컴파일, 실행 

코드 생성 및 빌드 

    ./build.sh 

서버 및 클라이언트 실행 

    cd rundriver 
     ./genwasmclient.sh&& go run server.go

    화면에 나오는 링크를 브라우저에서 실행
    open admin web
    http://localhost:24201/
    open client web
    http://localhost:24101/
    

# windows 에서 작동시키려면?

signalhandlewin을 사용하는 rundriver/serverwin.go 를 사용하시면 됩니다. 


# packet marshaler(serializer)

기본적으로 gob 를 사용하게 되어 있습니다. (gos_gob 패키지)

json을 사용하고 싶으면 gos_json 을 사용하도록 수정하고 

msgp(MessagePack) 를 사용하고 싶으면 gos_msgp를 사용하도록 수정하면 됩니다. 

그외 다른 marshaler를 사용하려면 위 세개의 패키지를 참고해서 marshal 패키지를 만드시면 되겠습니다. 

