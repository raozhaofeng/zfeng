import {api} from 'boot/axios';

// captcha 验证码
export const CaptchaAPI = () => {
    return api.get('/captcha/generate', {showLoading: false} as any);
};

// 管理员登陆
export const LoginAPI = (params: any) => {
    return api.post('/login', params);
};

// 首页信息
export const IndexAPI = () => {
    return api.post('/index')
}

// 管理资料
export const AdminInfoAPI = () => {
    return api.post('/info')
}

// 数据库表列表
export const TablesIndexAPI = () => {
    return api.post('/tables/index', {}, {showLoading: false} as any)
}

// 数据库表字段
export const TablesColumnsAPI = (params: any) => {
    return api.post('/tables/columns', params)
}

// 语言下载
export const DictionaryDownloadAPI = () => {
    return api.post('/dictionary/download')
}

// 语言上传
export const DictionaryUploadAPI = () => {
    return api.post('/dictionary/upload')
}