import {tableTemplate} from 'src/utils/template/table';

// 生成数据表格
export const generateTableFunc = (conf: any): string => {
    let content = tableTemplate;

    content = content.replace(/\{\{updateKey}}/g, conf.basic.updateKey)
    content = content.replace(/\{\{confBasic}}/g, JSON.stringify(conf.basic).replace(/"(\w+)":/g, '$1:').replace(/"/g, '\''))
    content = content.replace(/\{\{confSearch}}/g, JSON.stringify(conf.search).replace(/"(\w+)":/g, '$1:').replace(/"/g, '\''))
    content = content.replace(/\{\{confTools}}/g, JSON.stringify(conf.tools).replace(/"(\w+)":/g, '$1:').replace(/"/g, '\''))

    // 处理columns
    const tableColumns = [] as any
    const visibleColumns = [] as any
    for (let i = 0; i < conf.columns.length; i++) {
        const column = conf.columns[i]
        tableColumns.push({
            label: column.label,
            field: column.field,
            name: column.field,
            type: column.type,
            align: 'left',
            sortable: true,
            data: column.data
        })
        visibleColumns.push(column.field)
    }
    // 新增操作栏
    if (conf.columnOptions.length > 0) {
        tableColumns.push({
            label: '操作',
            field: '',
            name: 'options',
            type: 'options',
            align: 'left',
            items: conf.columnOptions
        })
        visibleColumns.push('options')
    }
    content = content.replace(/\{\{confColumns}}/g, JSON.stringify(tableColumns).replace(/"(\w+)":/g, '$1:').replace(/"/g, '\''))
    content = content.replace(/\{\{confVisibleColumns}}/g, JSON.stringify(visibleColumns).replace(/"/g, '\''))
    content = content.replace(/\{\{visibleStorage}}/g, '_table_visible_' + Math.ceil(Math.random()*1000000))

    return content
}