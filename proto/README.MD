```
protoc --go_out=. --go_opt=paths=source_relative proto/actionpb/action.proto

```

```
protoc --go_out=. --go_opt=paths=source_relative proto/global.proto
protoc --go_out=. --go_opt=paths=source_relative proto/transformpb/transform.proto
protoc --go_out=. --go_opt=paths=source_relative proto/objectpb/object.proto
protoc --go_out=. --go_opt=paths=source_relative proto/messagepb/message.proto 
protoc --go_out=. --go_opt=paths=source_relative proto/interactpb/interact.proto 
protoc --go_out=. --go_opt=paths=source_relative proto/pingpb/ping.proto 
protoc --go_out=. --go_opt=paths=source_relative proto/soundpb/sound.proto 
protoc --go_out=. --go_opt=paths=source_relative proto/animationpb/animation.proto 
```