<template>
    <div>
        <q-list>
            <q-item v-for="(conversation, conversationIndex) in conversationList" :key="conversationIndex" clickable
                    @click="showChatDialogFunc(conversation)">
                <q-item-section avatar>
                    <q-avatar size="50px" class="q-my-xs">
                        <q-img no-spinner :src="imageSrc(conversation.user_info.avatar)" width="60px">
                            <div class="absolute-full flex flex-center bg-grey" style="opacity: 0.9"
                                 v-if="!conversation.user_info.online"></div>
                        </q-img>
                    </q-avatar>
                </q-item-section>
                <q-item-section>
                    <div class="full-width">
                        <div class="text-body2 ellipsis"
                             :class="conversation.user_info.online ? 'text-black' : 'text-grey'">
                            {{ conversation.user_info.username }}
                        </div>
                        <div class="text-caption text-grey ellipsis" v-if="conversation.data == null">还未发送消息</div>
                        <div class="text-caption text-grey ellipsis"
                             v-else-if="conversation.data.type === WebsocketMessageTypeText">
                            {{ conversation.data.data }}
                        </div>
                        <div class="text-caption text-grey ellipsis" v-else>[图片]</div>
                        <div class="text-caption ellipsis"
                             :class="conversation.user_info.online ? 'text-secondary' : 'text-grey'">
                            {{ conversation.user_info.address }}
                        </div>
                    </div>
                </q-item-section>
                <q-item-section side>
                    <div class="full-height column justify-start">
                        <div class="text-caption text-grey">
                            {{ formatDate(conversation.updated_at, 'HH:MM') }}
                        </div>
                        <div class="q-mt-sm">
                            <q-badge rounded color="red" :label="conversation.unread" v-if="conversation.unread > 0"/>
                        </div>
                    </div>
                </q-item-section>
            </q-item>
        </q-list>

        <q-dialog v-model="chatDialog" @hide="currentConversation = null">
            <q-card class="full-width">
                <q-card-section class="bg-grey-2">
                    <div class="row justify-start q-gutter-sm items-center">
                        <div>
                            <q-avatar size="50px" class="q-my-xs">
                                <q-img no-spinner :src="imageSrc(currentConversation.user_info.avatar)" width="60px">
                                    <div class="absolute-full flex flex-center bg-grey" style="opacity: 0.9"
                                         v-if="!currentConversation.user_info.online"></div>
                                </q-img>
                            </q-avatar>
                        </div>
                        <div>
                            <div class="text-body2 ellipsis"
                                 :class="currentConversation.user_info.online ? 'text-black' : 'text-grey'">
                                {{ currentConversation.user_info.username }}
                            </div>
                            <div class="text-caption ellipsis"
                                 :class="currentConversation.user_info.online ? 'text-secondary' : 'text-grey'">
                                {{ currentConversation.user_info.address }}
                            </div>
                        </div>
                    </div>
                </q-card-section>
                <q-separator/>

                <q-card-section class="scroll bg-grey-2" style="padding: 0 0 0 8px">
                    <q-scroll-area ref="chatMessageRef" style="height: 50vh;"
                                   :thumb-style="{right: '2px', borderRadius: '5px', width: '5px'}"
                                   content-active-style="padding-right: 10px" content-style="padding-right: 10px">
                        <div v-for="(message, messageIndex) in chatMessageList" :key="messageIndex">
                            <q-chat-message :avatar="imageSrc(adminInfo.avatar)"
                                            v-if="message.role_id === WebsocketMessageRoleAdminToUser || message.role_id === WebsocketMessageRoleAdminToTourist"
                                            :text-html="message.type === WebsocketMessageTypeImage"
                                            :text='[message.type === WebsocketMessageTypeImage ? "<img src=\"" + imageSrc(message.data) + "\" width=\"100%\" />" : message.data]'
                                            :stamp="agoFormatDate(message.created_at)" sent
                                            bg-color="primary" text-color="white"/>
                            <q-chat-message :avatar="imageSrc(currentConversation.user_info.avatar)" v-else
                                            :text-html="message.type === WebsocketMessageTypeImage"
                                            :text='[message.type === WebsocketMessageTypeImage ? "<img src=\"" + imageSrc(message.data) + "\" width=\"100%\" />" : message.data]'
                                            :stamp="agoFormatDate(message.created_at)"
                                            bg-color="white"/>
                        </div>
                    </q-scroll-area>
                </q-card-section>
                <q-separator/>

                <!--                发送内容消息框-->
                <q-card-actions class="no-padding bg-grey-2">
                    <div class="row">
                        <q-uploader flat auto-upload style="width: 36px; height: 36px; line-height: 36px"
                                    class="chat-tool-uploader"
                                    :url='baseURL + "/upload"' field-name='file'
                                    @uploaded='sendImageFunc'
                                    :headers="[{name: 'Token', value: $store.state.user.token}, {name: 'Token-Key', value: $store.state.user.tokenKey}]">
                            <template v-slot:header></template>
                            <template v-slot:list='scope'>
                                <div @click='scope.pickFiles' class="bg-grey-2 no-padding full-height"
                                     style="overflow: hidden">
                                    <q-uploader-add-trigger/>
                                    <q-btn dense flat square color="primary" icon="sym_o_image"
                                           style="min-width: 0"/>
                                </div>
                            </template>
                        </q-uploader>
                    </div>
                    <div class="full-width">
                        <q-input v-model="msgText" autofocus type="textarea" input-style="padding: 4px 16px 16px 16px"
                                 placeholder="输入内容..."
                                 @keyup.enter="sendTextFunc"></q-input>
                    </div>
                </q-card-actions>
            </q-card>
        </q-dialog>
    </div>
</template>

<script lang="ts">
import store from 'src/store';
import {Notify} from 'quasar';
import {onMounted, reactive, toRefs, ref, watch} from 'vue'
import {imageSrc} from 'src/utils';
import {agoFormatDate, formatDate} from 'src/utils/formatDate';
import {
    websocketConn,
    WebsocketMessageTypeOnline,
    WebsocketMessageTypeOffline,
    WebsocketMessageRoleAdminToUser,
    WebsocketMessageRoleAdminToTourist,
    ChatSessionTypeAdminToTourist,
    WebsocketMessageTypeText,
    WebsocketMessageTypeImage,
} from 'boot/websocket';
import {
    SendMessageAPI,
    ConversationMessageAPI,
    ConversationIndexAPI,
    ClearUnreadMessageAPI,
    ConversationUserInfoAPI
} from 'src/api/chat';

export default {
    name: 'ChatConversation',
    props: {
        drawer: {type: Boolean, default: false}
    },
    setup(props: any) {
        const chatMessageRef = ref(null) as any
        const state = reactive({
            baseURL: process.env.baseURL,
            chatDialog: false,
            adminInfo: JSON.parse(JSON.stringify(store.state.user.info)),
            currentConversation: null as any,
            msgText: '',
            chatMessageList: [] as any,
            conversationList: [] as any
        })

        onMounted(() => {
            //  请求聊天会话列表
            ConversationIndexAPI().then((res: any) => {
                state.conversationList = res
            })

            //  消息通知
            if (websocketConn != null) {
                websocketConn.messageEvent = (msg: any) => {
                    const userInfo = JSON.parse(JSON.stringify(store.state.user.info))
                    switch (msg.type) {
                        // 上线通知
                        case WebsocketMessageTypeOnline:
                            userInfo.online_nums++
                            store.commit('user/updateInfo', userInfo)
                            operateConversation(msg.session_id, 'status', true)
                            break
                        // 下线通知
                        case WebsocketMessageTypeOffline:
                            if (userInfo.online_nums > 0) {
                                userInfo.online_nums--
                            }
                            store.commit('user/updateInfo', userInfo)
                            operateConversation(msg.session_id, 'status', false)
                            break
                        // 聊天消息
                        case WebsocketMessageTypeText:
                        case WebsocketMessageTypeImage:
                            //  发出消息声音
                            const audio = new Audio(imageSrc('/assets/mp3/msg.mp3'))
                            audio.play().then()
                            //  如果是当前对话框,那么添加消息， 如果不是, 那么追加未读消息
                            if (state.currentConversation != null) {
                                state.chatMessageList.push(msg)
                                ClearUnreadMessageAPI({id: state.currentConversation.id}).then()
                            } else {
                                //  未读消息 +1, 并且排序第一位
                                console.log(msg)
                                operateConversation(msg.session_id, 'unread', msg)
                            }
                            break
                    }
                }
            }
        })

        // 获取操作会话
        const operateConversation = (conversationId: number, opt: string, val: any) => {
            let addConversation = true
            for (let i = 0; i < state.conversationList.length; i++) {
                if (state.conversationList[i].id === conversationId) {
                    switch (opt) {
                        case 'status':
                            state.conversationList[i].user_info.online = val
                            break
                        case 'unread':
                            //  排序到第一位, 并且发送提示消息
                            state.conversationList[i].unread++
                            state.conversationList[i].data = {type: val.type, data: val.data}

                            //  排序到第一位
                            const conversationTmp = state.conversationList[i]
                            state.conversationList.splice(i, 1)
                            state.conversationList.unshift(conversationTmp)
                            addConversation = false
                            break
                    }
                    break
                }
            }

            //  新增会话对象
            if (opt === 'unread' && addConversation) {
                addConversationFunc(val)
            }

            //  如果是新消息, 并且没有打开聊天, 那么提示消息
            if (opt === 'unread' && !props.drawer) {
                const userInfo = JSON.parse(JSON.stringify(store.state.user.info))
                userInfo.unread_nums++
                store.commit('user/updateInfo', userInfo)
                const notifyMessage = val.type === WebsocketMessageTypeText ? '您有新的聊天消息【' + val.data + '】' : '您有新的聊天消息【图片】'
                Notify.create({type: 'positive', position: 'top-right', timeout: 3000, message: notifyMessage})
            }
        }

        // 新增会话
        const addConversationFunc = (msg: any) => {
            ConversationUserInfoAPI({id: msg.session_id}).then((res: any) => {
                state.conversationList.unshift(res)
            })
        }

        // 显示聊天框
        const showChatDialogFunc = (conversation: any) => {
            state.currentConversation = conversation
            //  请求聊天记录
            ConversationMessageAPI({id: conversation.id}).then((res: any) => {
                state.chatMessageList = res
            })
            //  删除未读消息
            ClearUnreadMessageAPI({id: conversation.id}).then(() => {
                conversation.unread = 0
            })
            state.chatDialog = true
        }

        // 发送文本消息
        const sendTextFunc = () => {
            if (state.msgText !== '' && state.msgText !== '\n') {
                sendMessageFunc(WebsocketMessageTypeText, state.msgText)
            }
            state.msgText = ''
        }

        // 发送图片消息
        const sendImageFunc = (info: any) => {
            const fileURL = JSON.parse(info.xhr.response).data;
            sendMessageFunc(WebsocketMessageTypeImage, fileURL)
        }

        // 发送消息
        const sendMessageFunc = (msgType: number, msgContent: string) => {
            // 加载消息列表
            state.chatMessageList.push({
                role_id: state.currentConversation.type === ChatSessionTypeAdminToTourist ? WebsocketMessageRoleAdminToTourist : WebsocketMessageRoleAdminToUser,
                sender_id: state.adminInfo.id,
                receiver_id: state.currentConversation.user_info.id,
                type: msgType,
                data: msgContent,
                extra: '',
                created_at: Date.parse(new Date().toString()) / 1000
            })

            // 发送消息
            SendMessageAPI({
                id: state.currentConversation.id,
                type: msgType,
                message: msgContent,
            }).then(() => {
                state.currentConversation.data = {type: msgType, data: msgContent}
            })
        }

        // 监听数据, 滚动到底部
        watch(() => state.chatMessageList, () => {
            setTimeout(() => {
                chatMessageRef.value.setScrollPosition('vertical', chatMessageRef.value.getScroll().verticalSize)
            }, 200)
        }, {deep: true})

        return {
            chatMessageRef,
            WebsocketMessageRoleAdminToUser,
            WebsocketMessageRoleAdminToTourist,
            WebsocketMessageTypeText,
            WebsocketMessageTypeImage,
            imageSrc,
            formatDate,
            agoFormatDate,
            showChatDialogFunc,
            sendTextFunc,
            sendImageFunc,
            ...toRefs(state)
        }
    }
}
</script>

<style scoped>

</style>
