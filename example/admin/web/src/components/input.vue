<template>
    <div v-if='type === inputTypeObject.select'>
        <q-select dense outlined v-model='value' :label='label' :readonly="readonly"
                  @update:model-value="$emit('update:modelValue', value)"
                  :options='[{label:"全部", value: ""}].concat(options)' style='min-width: 120px'
                  emit-value map-options></q-select>
    </div>
    <div v-else-if="type === inputTypeObject.selectNumber">
        <q-select dense outlined v-model.number='value' :label='label' :readonly="readonly"
                  @update:model-value="$emit('update:modelValue', value)"
                  :options='[{label:"全部", value: 0}].concat(options)' style='min-width: 120px'
                  emit-value map-options></q-select>
    </div>
    <div v-else-if='type === inputTypeObject.dateRangePicker'>
        <q-input dense outlined :label='label' v-model='value' :readonly="readonly"
                 :model-value='value !== "" && value !== undefined && value !== null ? value.from + " - " + value.to : ""'>
            <template v-slot:append>
                <q-icon name='event' class='cursor-pointer'>
                    <q-popup-proxy ref='qDateProxy' cover transition-show='scale'
                                   transition-hide='scale'>
                        <q-date range v-model='value' @update:model-value="$emit('update:modelValue', value)">
                            <div class='row items-center justify-end'>
                                <q-btn v-close-popup label='关闭' color='primary' flat/>
                            </div>
                        </q-date>
                    </q-popup-proxy>
                </q-icon>
            </template>
        </q-input>
    </div>
    <div v-else-if='type === inputTypeObject.datePicker'>
        <q-input dense outlined :label='label' v-model='value' :readonly="readonly"
                 :model-value='typeof value === "number" ? date.formatDate(value*1000, "YYYY/MM/DD") : value'>
            <template v-slot:append>
                <q-icon name='event' class='cursor-pointer'>
                    <q-popup-proxy ref='qDateProxy' cover transition-show='scale'
                                   transition-hide='scale'>
                        <q-date v-model='value' @update:model-value="$emit('update:modelValue', value)">
                            <div class='row items-center justify-end'>
                                <q-btn v-close-popup label='关闭' color='primary' flat/>
                            </div>
                        </q-date>
                    </q-popup-proxy>
                </q-icon>
            </template>
        </q-input>
    </div>
    <div v-else-if='type === inputTypeObject.toggle'>
        <q-toggle v-model='value' :label='label' :readonly="readonly"
                  @update:model-value="$emit('update:modelValue', value)"
                  :true-value='parseInt(options[0].value)'
                  :false-value='parseInt(options[1].value)'></q-toggle>
    </div>
    <div v-else-if='type === inputTypeObject.slider'>
        <q-slider v-model='value' label-always :readonly="readonly"
                  label :label-value='label + ": " + value'
                  @update:model-value="$emit('update:modelValue', value)"
                  :min='parseInt(options[0].value)' :max='parseInt(options[1].value)'></q-slider>
    </div>
    <div v-else-if='type === inputTypeObject.range' class='q-mt-xl'>
        <q-range :min='parseInt(options[0].value)' :max='parseInt(options[1].value)'
                 v-model='value' label-always :readonly="readonly"
                 @update:model-value="$emit('update:modelValue', value)"
                 :left-label-value="label + '【最小值: ' + value.min + '】'"
                 :right-label-value="label + '【最大值: ' + value.max + '】'"
        ></q-range>
    </div>

    <div v-else-if='type === inputTypeObject.timePicker'>
        <q-input dense outlined :label='label' mask='fulltime' :readonly="readonly"
                 @update:model-value="$emit('update:modelValue', value)"
                 v-model='value'>
            <template v-slot:append>
                <q-icon name='access_time' class='cursor-pointer'>
                    <q-popup-proxy cover transition-show='scale' transition-hide='scale'>
                        <q-time with-seconds v-model='value'>
                            <div class='row items-center justify-end'>
                                <q-btn v-close-popup label='关闭' color='primary' flat/>
                            </div>
                        </q-time>
                    </q-popup-proxy>
                </q-icon>
            </template>
        </q-input>
    </div>
    <div v-else-if='type === inputTypeObject.textarea'>
        <q-input dense outlined rows='5' :readonly="readonly"
                 v-model='value' :label='label'
                 @update:model-value="$emit('update:modelValue', value)"
                 type='textarea'></q-input>
    </div>
    <div v-else-if='type === inputTypeObject.editor'>
        <q-editor min-height='5rem' :readonly="readonly" ref="editorRef" :placeholder="label"
                  :toolbar="[
                      ['token'],
                      [
                          {label: $q.lang.editor.align,icon: $q.iconSet.editor.align,fixedLabel: true,list: 'only-icons',options: ['left', 'center', 'right', 'justify']},
                          {label: $q.lang.editor.align,icon: $q.iconSet.editor.align,fixedLabel: true,options: ['left', 'center', 'right', 'justify']}
                      ],
                      ['bold', 'italic', 'strike', 'underline', 'subscript', 'superscript'],
                      ['hr', 'link', 'custom_btn'],
                      ['upload'],
                      ['print', 'fullscreen'],
                      [
                          {label: $q.lang.editor.formatting,icon: $q.iconSet.editor.formatting,list: 'no-icons',options: ['p','h1','h2','h3','h4','h5','h6','code']},
                          {label: $q.lang.editor.fontSize,icon: $q.iconSet.editor.fontSize,fixedLabel: true,fixedIcon: true,list: 'no-icons',options: ['size-1','size-2','size-3','size-4','size-5','size-6','size-7']},
                          {label: $q.lang.editor.defaultFont,icon: $q.iconSet.editor.font,fixedIcon: true,list: 'no-icons',options: ['default_font','arial','arial_black','comic_sans','courier_new','impact','lucida_grande','times_new_roman','verdana']},
                          'removeFormat'
                      ],
                      ['quote', 'unordered', 'ordered', 'outdent', 'indent'],
                      ['undo', 'redo'],
                      ['viewsource']
                  ]"
                  :fonts="{arial: 'Arial',arial_black: 'Arial Black',comic_sans: 'Comic Sans MS',courier_new: 'Courier New',impact: 'Impact',lucida_grande: 'Lucida Grande',times_new_roman: 'Times New Roman',verdana: 'Verdana'}"
                  v-model='value'
                  @update:model-value="$emit('update:modelValue', value)" :dense='$q.screen.lt.md'
        >
            <template v-slot:token>
                <q-btn-dropdown dense no-caps ref="editorTokenRef" no-wrap unelevated color="white" text-color="primary"
                                label="图片｜颜色" size="sm">
                    <q-list dense>
                        <q-item class="q-mb-md">
                            <q-uploader flat auto-upload style="height: 36px"
                                        @uploaded="editorUploadedEventFunc"
                                        :headers="[{name: 'Token', value: $store.state.user.token}, {name: 'Token-Key', value: $store.state.user.tokenKey}]"
                                        :url='baseURL + "/upload"' field-name='file'>
                                <template v-slot:header></template>
                                <template v-slot:list='scope'>
                                    <div @click='scope.pickFiles'>
                                        <q-uploader-add-trigger/>
                                        <div class="text-body2 text-grey" style="height: 36px; line-height: 36px;">
                                            选择需要上传图片...
                                        </div>
                                    </div>
                                </template>
                            </q-uploader>
                        </q-item>
                        <q-item tag="label" clickable
                                @click="editorEditTextColorFunc('backColor', editorBackgroundColor)">
                            <q-item-section side>
                                <q-icon name="highlight"/>
                            </q-item-section>
                            <q-item-section>
                                <q-color
                                    v-model="editorBackgroundColor"
                                    default-view="palette"
                                    no-header
                                    no-footer
                                    :palette="['#ffccccaa', '#ffe6ccaa', '#ffffccaa', '#ccffccaa', '#ccffe6aa', '#ccffffaa', '#cce6ffaa', '#ccccffaa', '#e6ccffaa', '#ffccffaa', '#ff0000aa', '#ff8000aa', '#ffff00aa', '#00ff00aa', '#00ff80aa', '#00ffffaa', '#0080ffaa', '#0000ffaa', '#8000ffaa', '#ff00ffaa', '#ff00ffaa']"
                                />
                            </q-item-section>
                        </q-item>

                        <q-item tag="label" clickable
                                @click="editorEditTextColorFunc('foreColor', editorTextColor)">
                            <q-item-section side>
                                <q-icon name="format_paint"/>
                            </q-item-section>
                            <q-item-section>
                                <q-color
                                    v-model="editorTextColor"
                                    default-view="palette"
                                    no-header
                                    no-footer
                                    :palette="['#000000', '#ff0000', '#ff8000', '#ffff00', '#00ff00', '#00ffff', '#0080ff', '#0000ff', '#8000ff', '#ff00ff']"
                                />
                            </q-item-section>
                        </q-item>
                    </q-list>
                </q-btn-dropdown>

            </template>
        </q-editor>

    </div>
    <div v-else-if='type === inputTypeObject.file'>
        <q-uploader auto-upload :readonly="readonly"
                    :headers="[{name: 'Token', value: $store.state.user.token}, {name: 'Token-Key', value: $store.state.user.tokenKey}]"
                    @uploading="uploadStartFunc"
                    @failed="uploadFailedFunc"
                    @uploaded="uploadedEventFunc"
                    :url='baseURL + "/upload"' field-name='file' :label='label'>
            <template v-slot:list>
                <div>{{ value }}</div>
            </template>
        </q-uploader>
    </div>
    <div v-else-if='type === inputTypeObject.image || type === inputTypeObject.images'>
        <q-uploader auto-upload :readonly="readonly"
                    @uploading="uploadStartFunc"
                    @failed="uploadFailedFunc"
                    @uploaded='uploadedEventFunc' :headers="[{name: 'Token', value: $store.state.user.token}, {name: 'Token-Key', value: $store.state.user.tokenKey}]"
                    :url='baseURL + "/upload"' field-name='file' :label='label'>
            <template v-slot:list>
                <q-list separator>
                    <q-item v-for='(imageItem, imageIndex) in value' :key='imageIndex'>
                        <q-item-section>
                            <q-item-label class='full-width ellipsis'>
                                {{ imageItem.label }}
                            </q-item-label>
                        </q-item-section>
                        <q-item-section class='gt-xs' thumbnail>
                            <img
                                :src='imageSrc(imageItem.value)'
                                alt=''/>
                        </q-item-section>
                        <q-item-section side>
                            <q-btn class='gt-xs' size='12px' flat dense round icon='delete'
                                   @click='deleteUploadedEventFunc(imageIndex)'/>
                        </q-item-section>
                    </q-item>
                </q-list>
            </template>
        </q-uploader>
    </div>
    <div v-else-if="type === inputTypeObject.icon">
        <q-uploader flat auto-upload style="width: 50px"
                    :url='baseURL + "/upload"' field-name='file'
                    @uploading="uploadStartFunc"
                    @failed="uploadFailedFunc"
                    @uploaded='uploadedEventFunc'
                    :headers="[{name: 'Token', value: $store.state.user.token}, {name: 'Token-Key', value: $store.state.user.tokenKey}]">
            <template v-slot:header></template>
            <template v-slot:list='scope'>
                <div @click='scope.pickFiles' class="row justify-center no-padding">
                    <q-uploader-add-trigger/>
                    <q-avatar v-if="value.length === 0" size="50px">
                        <q-icon name="sym_o_add_photo_alternate" color="grey" size="50px"></q-icon>
                    </q-avatar>
                    <q-avatar v-else size="50px">
                        <q-img no-spinner :src="imageSrc(value[0].value)"></q-img>
                    </q-avatar>
                </div>
            </template>
        </q-uploader>
    </div>
    <div v-else-if='type === inputTypeObject.checkbox' class='q-gutter-sm'>
        <div class='q-ml-sm text-body2 text-bold text-grey'>{{ label }}</div>
        <q-checkbox v-for='(checkbox, checkboxIndex) in options' :readonly="readonly"
                    v-model='value[checkbox.value]' :label='checkbox.label'
                    class='no-margin'
                    @update:model-value="$emit('update:modelValue', value)"
                    :key='checkboxIndex'>
        </q-checkbox>
    </div>
    <div v-else-if="type === inputTypeObject.children" class='q-gutter-sm'>
        <div class="row justify-end q-mb-md q-gutter-xs">
            <div>
                <q-btn dense flat outline class="bg-secondary text-white q-pa-xs" size="xs"
                       @click="addChildrenFunc()"
                       label="添加"></q-btn>
            </div>
            <div>
                <q-btn dense flat outline class="bg-red text-white q-pa-xs" size="xs"
                       @click="delChildrenFunc()"
                       label="删除"></q-btn>
            </div>
        </div>
        <div v-for="(childrenValue, childrenValueIndex) in value"
             :key="childrenValueIndex">
            <div class="row q-gutter-md items-center">
                <div class="col q-mb-md"
                     v-for="(childrenData, childrenDataIndex) in data"
                     :key="childrenDataIndex">
                    <dynamic-input :type="childrenData.type" :label="childrenData.label"
                                   v-model="value[childrenValueIndex][childrenData.field]"
                                   :data="[]"></dynamic-input>
                </div>
            </div>
        </div>
    </div>
    <div v-else-if="type === inputTypeObject.number">
        <q-input dense outlined v-model.number='value' :label='label' type="number" :readonly="readonly"
                 @change="$emit('update:modelValue', value)"></q-input>
    </div>
    <div v-else-if="type === inputTypeObject.password">
        <q-input dense outlined v-model='value' :label='label' type="password" :readonly="readonly"
                 @change="$emit('update:modelValue', value)"></q-input>
    </div>
    <div v-else-if="type === inputTypeObject.color">
        <q-input dense outlined v-model="value" :label="label" @change="$emit('update:modelValue', value)">
            <template v-slot:append>
                <q-icon name="colorize" class="cursor-pointer">
                    <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                        <q-color v-model="value" @change="$emit('update:modelValue', value)"/>
                    </q-popup-proxy>
                </q-icon>
            </template>
        </q-input>
    </div>
    <div v-else>
        <q-input dense outlined v-model='value' :label='label' :readonly="readonly"
                 @change="$emit('update:modelValue', value)"></q-input>
    </div>
</template>

<script lang='ts'>
import {date, Loading, QSpinnerBars} from 'quasar';
import {imageSrc, negativeNotify} from 'src/utils';
import {inputTypeObject} from 'src/utils/define';
import {reactive, toRefs, watch, ref} from 'vue';

export default {
    name: 'DynamicInput',
    props: {
        modelValue: {
            type: undefined,
            default: () => {
                return '';
            }
        },
        label: {
            type: String
        },
        readonly: {
            type: Boolean
        },
        type: {
            type: String
        },
        data: {
            type: Array
        }
    },
    setup(props: any, context: any) {
        const editorRef = ref(null) as any
        const editorTokenRef = ref(null) as any
        const state = reactive({
            editorTextColor: '',
            editorBackgroundColor: '',
            baseURL: process.env.baseURL,
            value: props.modelValue,
            options: props.data
        });

        watch(props, (newProps) => {
            state.options = newProps.data;
        });

        // 预先处理多选框
        if (props.type === inputTypeObject.checkbox) {
            state.value = Object.assign({}, state.value);
            state.options.forEach((item: any) => {
                if (!state.value.hasOwnProperty(item.value)) {
                    state.value[item.value] = false;
                }
            });
        }

        //  预先处理子级框
        if (props.type === inputTypeObject.children) {
            if (state.value == null) {
                state.value = [] as any
            }
        }

        // 多选框处理
        if (props.type === inputTypeObject.selectNumber) {
            state.options.forEach((item: any) => {
                item.value = parseInt(item.value)
            })
        }

        // 预先处理图片
        if (props.type === inputTypeObject.image || props.type === inputTypeObject.icon || props.type === inputTypeObject.images) {
            if (state.value === undefined || state.value === '') {
                state.value = [] as any;
            }

            if ((props.type === inputTypeObject.icon || props.type === inputTypeObject.image) && state.value !== '' && Object.keys(state.value).length !== 0) {
                state.value = [{label: '图片简介说明', value: state.value}];
            }
        }

        // 处理富文本
        if (props.type === inputTypeObject.editor) {
            // 默认显示源码
            setTimeout(() => {
                editorRef.value.runCmd('viewsource', true)
            }, 100)

            if (state.value === undefined) {
                state.value = '';
            }
        }

        // 上传文件事件
        const uploadedEventFunc = (info: any) => {
            const fileURL = JSON.parse(info.xhr.response).data;
            if (props.type === inputTypeObject.image || props.type === inputTypeObject.icon) {
                state.value = [{label: '图片简介说明', value: fileURL}];
                context.emit('update:modelValue', fileURL);
            } else if (props.type === inputTypeObject.file) {
                state.value = fileURL
                context.emit('update:modelValue', state.value);
            } else {
                state.value.push({label: '图片简介说明', value: fileURL});
                context.emit('update:modelValue', state.value);
            }
            Loading.hide()
        };

        // 开始上传方法
        const uploadStartFunc = () => {
            Loading.show({
                spinner: QSpinnerBars,
                spinnerColor: 'secondary',
                spinnerSize: 50,
                message: '开始上传文件....'
            });
        }

        // 上传失败方法
        const uploadFailedFunc = () => {
            negativeNotify('文件上传失败, 请检查文件是否符合邀请...')
            Loading.hide()
        }

        // 删除上传的图片
        const deleteUploadedEventFunc = (index: number) => {
            state.value.splice(index, 1);
            if (props.type === inputTypeObject.image || props.type === inputTypeObject.icon) {
                context.emit('update:modelValue', '');
            } else {
                context.emit('update:modelValue', state.value);
            }
        };

        // 富文本编辑器上传图片
        const editorUploadedEventFunc = (info: any) => {
            const fileURL = JSON.parse(info.xhr.response).data;
            state.value += '<img src="' + imageSrc(fileURL) + '" />'
            editorTokenRef.value.hide()
            Loading.hide()
        }

        // 编辑文本颜色
        const editorEditTextColorFunc = (cmd: string, color: string) => {
            editorTokenRef.value.hide()
            editorRef.value.runCmd(cmd, color)
            editorRef.value.focus()
        }

        // 新增子级
        const addChildrenFunc = () => {
            let oldData = props.data
            let tmpValue = {} as any

            for (let i = 0; i < oldData.length; i++) {
                tmpValue[oldData[i].field] = ''
            }
            state.value.push(tmpValue)
        }

        //  删除子级
        const delChildrenFunc = () => {
            state.value.pop()
        }

        return {
            date,
            inputTypeObject,
            imageSrc,
            uploadStartFunc,
            uploadFailedFunc,
            uploadedEventFunc,
            deleteUploadedEventFunc,
            editorRef,
            editorTokenRef,
            editorUploadedEventFunc,
            editorEditTextColorFunc,
            addChildrenFunc,
            delChildrenFunc,
            ...toRefs(state)
        };
    }
};
</script>

<style scoped>

</style>
