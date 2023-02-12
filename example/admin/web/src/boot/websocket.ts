import {boot} from 'quasar/wrappers';
import {imageSrc, negativeNotify, positiveNotify} from 'src/utils';
import store from 'src/store'

const WebsocketMessageTypeRegisterAudio = 1000      //  音频消息
export const ChatSessionTypeAdminToUser = 10        //  管理对用户
export const ChatSessionTypeAdminToTourist = 11     //  管理对临时用户
export const WebsocketMessageRoleAdminToUser = 12          //  管理发送用户
export const WebsocketMessageRoleAdminToTourist = 13        //  管理发送临时用户
export const WebsocketMessageTypeText = 1           //  文本消息
export const WebsocketMessageTypeImage = 2          //  图片消息
export const WebsocketMessageTypeOnline = 1010             //  上线通知
export const WebsocketMessageTypeOffline = 1011            //  下线通知

let websocketConn = null as any                     //  websocket对象

export const connectWebsocket = () => {
    if (websocketConn == null) {
        const protocol = document.location.protocol === 'https:' ? 'wss:' : 'ws:'
        websocketConn = new WebSocket(protocol + process.env.baseURL + '/ws?token=' + store.state.user.token + '&key=' + store.state.user.tokenKey)
        // 打开websocket消息
        websocketConn.onopen = () => {
            console.log('连接websocket成功...')
        }

        // 消息通知
        websocketConn.messageEvent = null as any

        // 发送websocket消息
        websocketConn.onmessage = (event: any) => {
            const message = JSON.parse(event.data)
            switch (message.type) {
                //  播放声音
                case WebsocketMessageTypeRegisterAudio:
                    playSound(message)
                    break;
                default:
                    //  聊天消息
                    if (websocketConn.messageEvent != null) {
                        websocketConn.messageEvent(message)
                    }
            }
        }

        // 关闭websocket消息
        websocketConn.onclose = () => {
            websocketConn = null
            console.log('websocket已关闭')
        }
    }
}

// 播放音频消息
const playSound = (message: any) => {
    const audio = new Audio(imageSrc(message.extra))
    audio.play().then(() => {
        positiveNotify(message.message)
    }).catch(() => {
        negativeNotify('音频播放失败, 请点击页面激活互动....')
    })
}

// 连接websocket
connectWebsocket()

export {websocketConn}
export default boot(() => {
    //  启动文件
})
