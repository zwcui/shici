package controllers

import (
	"net/http"
	"golang.org/x/net/websocket"
	"strconv"
	"encoding/json"
	"sort"
	"crypto/md5"
	"encoding/hex"
	"strings"
	"shici/models"
	"shici/util"
	"flag"
)

//socket连接池
var users map[int64]models.SocketConnection

type WSServer struct {
	ListenAddr string
}

func init(){
	go startSocket()
}

func startSocket(){
	util.Logger.Info("--------websocket--------start--------")
	addr := flag.String("a", ":6666", "websocket server listen address")
	flag.Parse()
	wsServer := &WSServer{
		ListenAddr : *addr,
	}
	wsServer.Start()
	util.Logger.Info("--------websocket--------end--------")
}

func (this *WSServer) Handler (conn *websocket.Conn) {
	if users == nil {
		users = make(map[int64]models.SocketConnection)
	}

	util.Logger.Info("a new ws conn: conn.RemoteAddr()="+conn.RemoteAddr().String()+"  conn.LocalAddr()="+conn.LocalAddr().String())
	var err error
	for {
		var reply string
		var socketMessage models.SocketMessage
		err = websocket.Message.Receive(conn, &reply)
		if err != nil {
			for k, v := range users {
				if v.Conn == conn {
					delete(users, k)
				}
			}
			util.Logger.Info("receive err:",err.Error())
			break
		}

		//util.Logger.Info("-----reply----  "+reply)
		err = json.Unmarshal([]byte(reply), &socketMessage)
		if err != nil {
			util.Logger.Info("----socketMessage--json.Unmarshal--err---- "+err.Error())
			continue
		}

		//验签
		if socketMessage.MessageSign != signMessage(socketMessage) {
			util.Logger.Info("----socketMessage--签名验证失败")
			break
		}

		//建立心跳
		if socketMessage.MessageType == 0 {
			handleHeartbeat(&socketMessage, conn)
		}

		//消息聊天
		if socketMessage.MessageType == 1 {
			var userMessage models.UserSocketMessage
			err = json.Unmarshal([]byte(socketMessage.MessageContent), &userMessage)
			if err != nil {
				util.Logger.Info("----userMessage--json.Unmarshal--err---- "+err.Error())
				continue
			}
			handleUserMessage(&userMessage, conn, socketMessage.MessageExpireTime)
			util.Logger.Info("Received from client: " + userMessage.Content+"  from "+strconv.FormatInt(userMessage.FromUid, 10)+" to "+strconv.FormatInt(userMessage.ToUid, 10))
		}
	}
}

func (this *WSServer) Start() (error) {
	http.Handle("/ws", websocket.Handler(this.Handler))
	util.Logger.Info("websocket----begin to listen")
	err := http.ListenAndServe(this.ListenAddr, nil)
	if err != nil {
		util.Logger.Info("ListenAndServe:", err)
		return err
	}
	util.Logger.Info("websocket----start end")
	return nil
}

//处理心跳
func handleHeartbeat(socketMessage *models.SocketMessage, conn *websocket.Conn){
	if _, ok := users[socketMessage.MessageSendUid]; !ok {
		var socketConnection models.SocketConnection
		socketConnection.Conn = conn
		socketConnection.ExpireTime = socketMessage.MessageExpireTime
		socketConnection.Token = socketMessage.MessageToken
		socketConnection.AppNo = socketMessage.MessageAppNo
		users[socketMessage.MessageSendUid] = socketConnection
		util.Logger.Info("-----socketConnection.heartbeat.ExpireTime----start-"+strconv.FormatInt(socketMessage.MessageSendUid, 10)+"--"+strconv.FormatInt(users[socketMessage.MessageSendUid].ExpireTime, 10))
	} else {
		//以token和appNo作为唯一标示
		if users[socketMessage.MessageSendUid].Token != socketMessage.MessageToken || users[socketMessage.MessageSendUid].AppNo != socketMessage.MessageAppNo {
			//账户被挤下线
			util.Logger.Info("-------您的账户已在其他地方登陆-------"+strconv.FormatInt(socketMessage.MessageSendUid, 10))
			var alert models.AlertMessage
			alert.AlertCode = "WebSocket400"
			alert.AlertMessage = "您的账户已在其他地方登陆"
			alertJsonByte, err := json.Marshal(alert)
			if err != nil {
				util.Logger.Info("---json to string---您的账户已在其他地方登陆----err:"+err.Error())
				//return
			}

			var replySocketMessage models.SocketMessage
			replySocketMessage.MessageType = 2
			replySocketMessage.MessageSendTime = util.UnixOfBeijingTime()
			replySocketMessage.MessageSendUid = socketMessage.MessageSendUid
			replySocketMessage.MessageExpireTime = util.UnixOfBeijingTime()+3
			replySocketMessage.MessageContent = string(alertJsonByte)
			replySocketMessage.MessageToken = socketMessage.MessageToken
			replySocketMessage.MessageAppNo = socketMessage.MessageAppNo
			replySocketMessage.MessageSign = signMessage(replySocketMessage)

			replySocketMessageJsonByte, err := json.Marshal(replySocketMessage)
			if err != nil {
				util.Logger.Info("---json to string---replySocketMessage----err:"+err.Error())
				//return
			}

			if err := websocket.Message.Send(users[socketMessage.MessageSendUid].Conn, string(replySocketMessageJsonByte)); err != nil {
				util.Logger.Info("----userMessage--websocket.Message.Send 您的账户已在其他地方登陆 err:", err.Error())
				//移除出错的链接
				delete(users, socketMessage.MessageSendUid)
			}
			//新登陆的账户
			var socketConnection models.SocketConnection
			socketConnection.Conn = conn
			socketConnection.ExpireTime = socketMessage.MessageExpireTime
			socketConnection.Token = socketMessage.MessageToken
			socketConnection.AppNo = socketMessage.MessageAppNo
			users[socketMessage.MessageSendUid] = socketConnection
		} else {
			util.Logger.Info("-----socketConnection.heartbeat.ExpireTime-----"+strconv.FormatInt(socketMessage.MessageSendUid, 10)+"--"+strconv.FormatInt(users[socketMessage.MessageSendUid].ExpireTime, 10))
			socketConnection := users[socketMessage.MessageSendUid]
			socketConnection.ExpireTime = socketMessage.MessageExpireTime
			users[socketMessage.MessageSendUid] = socketConnection
		}
	}

	//查看redis缓存消息

}

//处理聊天
func handleUserMessage(userMessage *models.UserSocketMessage, conn *websocket.Conn, expireTime int64) {
	//聊天信息入库MongoDB

}

//签名
func signMessage(socketMessage models.SocketMessage) string {
	params := make(map[string]string)
	params["messageType"] = strconv.Itoa(socketMessage.MessageType)
	params["messageSendTime"] = strconv.FormatInt(socketMessage.MessageSendTime, 10)
	params["messageSendUid"] = strconv.FormatInt(socketMessage.MessageSendUid, 10)
	params["messageExpireTime"] = strconv.FormatInt(socketMessage.MessageExpireTime, 10)
	params["messageContent"] = socketMessage.MessageContent
	params["messageToken"] = socketMessage.MessageToken
	params["messageAppNo"] = strconv.Itoa(socketMessage.MessageAppNo)

	keys := make([]string, len(params))

	i := 0
	for k := range params {
		keys[i] = k
		i++
	}

	sort.Strings(keys)

	strTemp := ""
	for _, key := range keys {
		strTemp = strTemp + key + "=" + params[key] + "&"
	}
	strTemp += "key=" + models.SOCKET_MESSAGE_SIGN_KEY

	hasher := md5.New()
	hasher.Write([]byte(strTemp))
	md5Str := hex.EncodeToString(hasher.Sum(nil))

	return strings.ToUpper(md5Str)
}

//每分钟检查失效的socket连接
func checkSocketHeartbeat(){
	util.Logger.Info("定时任务，每分钟检查失效的socket连接")
	for uId, socketConnection := range users {
		util.Logger.Info("-----定时任务  遍历users-----  util.UnixOfBeijingTime()="+strconv.FormatInt(util.UnixOfBeijingTime(), 10)+"   uid="+strconv.FormatInt(uId, 10)+"   ExpireTime="+strconv.FormatInt(socketConnection.ExpireTime, 10))
		if (socketConnection.ExpireTime + 15) <= util.UnixOfBeijingTime() {
			delete(users, uId)
		}
	}
}


