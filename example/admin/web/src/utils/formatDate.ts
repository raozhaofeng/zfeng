import {date} from 'quasar'

// 格式化时间戳
export const formatDate = (timestamp: number, format: string): string => {
    return date.formatDate(timestamp * 1000, format)
}

// 时间戳过去时间
export const agoFormatDate = (timestamp: number): string => {
    timestamp = timestamp * 1000
    const minute = 1000 * 60;
    const hour = minute * 60;
    const day = hour * 24;
    const month = day * 30;
    const nowTime = new Date().getTime();
    const diffValue = nowTime - timestamp
    if (diffValue < 0) {
        return '不久前';
    }

    const monthDifference = diffValue / month;
    const weekDifference = diffValue / (7 * day);
    const dayDifference = diffValue / day;
    const hourDifference = diffValue / hour;
    const minDifference = diffValue / minute;

    if (monthDifference > 4) {
        return formatDate(timestamp, 'YYYY/MM/DD HH:ii:ss')
    } else if (monthDifference >= 1) {
        return parseInt(String(monthDifference)) + '月前';
    } else if (weekDifference >= 1) {
        return parseInt(String(weekDifference)) + '周前';
    } else if (dayDifference >= 1) {
        return parseInt(String(dayDifference)) + '天前';
    } else if (hourDifference >= 1) {
        return parseInt(String(hourDifference)) + '小时前';
    } else if (minDifference >= 1) {
        return parseInt(String(minDifference)) + '分钟前';
    }
    return '刚刚'
}