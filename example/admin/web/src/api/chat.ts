import {api} from 'boot/axios';

// 会话列表
export const ConversationIndexAPI = () => {
    return api.post('/chat/index', {}, {showLoading: false} as any)
}

// 会话消息
export const ConversationMessageAPI = (params: any) => {
    return api.post('/chat/message', params)
}

// 发送消息
export const SendMessageAPI = (params: any) => {
    return api.post('/chat/send', params, {showLoading: false} as any)
}

// 清除未读消息
export const ClearUnreadMessageAPI = (params: any) => {
    return api.post('/chat/unread', params, {showLoading: false} as any)
}

// 会话用户信息
export const ConversationUserInfoAPI = (params: any) => {
    return api.post('/chat/info', params, {showLoading: false} as any)
}