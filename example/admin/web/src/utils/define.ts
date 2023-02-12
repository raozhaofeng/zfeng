// 表单类型
export const inputTypeObject = {
    select: 'select',
    selectNumber: 'selectNumber',
    toggle: 'toggle',
    slider: 'slider',
    range: 'range',
    dateRangePicker: 'dateRangePicker',
    datePicker: 'datePicker',
    timePicker: 'timePicker',
    textarea: 'textarea',
    editor: 'editor',
    file: 'file',
    image: 'image',
    icon: 'icon',
    images: 'images',
    checkbox: 'checkbox',
    text: 'text',
    number: 'number',
    password: 'password',
    editText: 'editText',
    editTextarea: 'editTextarea',
    editNumber: 'editNumber',
    editToggle: 'editToggle',
    color: 'color',
    children: 'children',
    json: 'json',
    ip4: 'ip4'
};

// 默认表单配置
export const inputDefaultConf = {
    label: '',
    field: '',
    type: inputTypeObject.text,
}

// 默认表单options配置
export const inputDefaultOptionsConf = {label: '', value: ''}

// 颜色集合
export const quasarColorsObject = {
    primary: 'primary', secondary: 'secondary', accent: 'accent', dark: 'dark', positive: 'positive',
    negative: 'negative', info: 'info', warning: 'warning', red: 'red', pink: 'pink', purple: 'purple',
    deepPurple: 'deep-purple', indigo: 'indigo', blue: 'blue', lightBlue: 'light-blue', cyan: 'cyan',
    teal: 'teal', green: 'green', lightGreen: 'light-green', lime: 'lime', yellow: 'yellow', amber: 'amber',
    orange: 'orange', deepOrange: 'deep-orange', brown: 'brown', grey: 'grey', blueGrey: 'blue-grey'
}

// 默认按钮工具类型
export const editBtnTypeObject = {
    create: 'create',
    update: 'update',
    delete: 'delete',
    view: 'view',
    checkboxDelete: 'checkboxDelete',
}

// 默认按钮工具配置
export const editBtnDefaultConf = {
    label: '',
    url: '',
    color: '',
    type: editBtnTypeObject.create,
    tips: '',
    isShow: '',
    items: []
}

// 按钮工具判断是否显示
export const editBtnEvalIsShowFunc = (str: string, data: any): boolean => {
    if (str === '') {
        return true;
    }

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const {scope} = data

    // 需要添加， 不然会报错
    console.log(scope)
    return eval(str)
}