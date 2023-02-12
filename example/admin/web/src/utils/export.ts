import {exportFile} from 'quasar';
import {warningNotify} from 'src/utils/index';

// 下载csv文件
export const exportCSVFile = (oldColumns: any, visible: any, data: any) => {
    let columns = [] as any
    if (visible.length > 0) {
        for (let i = 0; i < oldColumns.length; i++) {
            if (visible.indexOf(oldColumns[i].field) > -1 && oldColumns[i].name !== 'options') {
                columns.push({label: oldColumns[i].label, field: oldColumns[i].field})
            }
        }
    } else {
        columns = oldColumns
    }

    const content = [columns.map((col: any) => col.label)].concat(
        data.map((row: any) => columns.map((col: any) => '"' + row[col.field] + '"').join(','))
    ).join('\r\n')

    const status = exportFile(
        'table-export.csv',
        content,
        'text/csv'
    )

    if (status !== true) {
        warningNotify('Browser denied file download...')
    }
}