import {Notify, Dialog} from 'quasar';

// 成功提示
export const positiveNotify = (msg: string) => {
    Notify.create({type: 'positive', position: 'top', timeout: 3000, message: msg});
};

// 失败提示
export const negativeNotify = (msg: string) => {
    Notify.create({type: 'negative', position: 'top', timeout: 3000, message: msg});
};

// 系统通知
export const warningNotify = (msg: string) => {
    Notify.create({type: 'warning', position: 'top-right', timeout: 3000, message: msg});
};

// 确认消息
export const confirmBoxDialog = (title: string, msg: string, onOk: any) => {
    Dialog.create({message: msg, title: title}).onOk(onOk);
};

// 图片方法
export const imageSrc = (url: string) => {
    if (url === '') {
        return '/icons/favicon.png';
    }
    if (url.indexOf('http') > -1 || url.indexOf('//') > -1) {
        return url
    }
    return rootDomain() + url
};

// 获取根域名
export const rootDomain = () => {
    return process.env.baseURL?.substring(0, process.env.baseURL.lastIndexOf('/'))
}