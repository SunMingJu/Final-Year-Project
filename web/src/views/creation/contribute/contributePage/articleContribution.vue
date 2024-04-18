<template>
    <div class="article-contribution">
        <!--Basic settings -->
        <div class="form-box animate__animated animate__bounceInRight" v-show="form.isShow">
            <h3> basic settings</h3>
            <el-form :model="form" ref="ruleFormRef" label-width="120px" label-position="left"
                :rules="rules.articleContributionRules">
                <el-form-item label="title" prop="title">
                    <el-input v-model="form.title" placeholder="Give the article a title~" />
                </el-form-item>
                <el-form-item class="form-item-middle" label="cover">
                    <el-upload class="cover-uploader" :action="uploadCoveration.action" :show-file-list="false"
                        :on-success="handleCover.handleFileSuccess" :on-error="handleCover.handleFileError"
                        :before-upload="handleCover.beforeFileUpload" :auto-upload="true"
                        :http-request="handleCover.RedefineUploadFile" accept=".png,.jpg,.jpeg">
                        <img v-if="uploadCoveration.FileUrl" :src="uploadCoveration.FileUrl" class="cover" />
                        <el-icon v-else class="cover-uploader-icon">
                            <Plus />
                        </el-icon>
                    </el-upload>
                </el-form-item>
                <el-form-item label="edit content">
                    <el-button v-removeFocus size="small" type="primary" :icon="Edit" round
                        @click="form.isShow = false">edit</el-button>
                </el-form-item>
                <el-form-item label="Release regularly" v-show="props.type != 'edit'">
                    <el-switch v-model="form.timing" />
                </el-form-item>
                <el-form-item label="selection period" v-show="form.timing" class="animate__animated animate__fadeIn">
                    <el-col :span="7">
                        <el-date-picker v-model="form.date1time" type="datetime" placeholder="Please select a scheduled release time" />
                    </el-col>
                </el-form-item>
                <el-form-item label="Enable comments">
                    <el-switch v-model="form.comments" />
                </el-form-item>
                <el-form-item label="Classification">
                    <el-cascader v-model="cascaderValue" placeholder="please select a type" :props="{
                        value: 'id',
                        label: 'label',
                        children: 'children'
                    }" :options="cascader" :show-all-levels="false" @change="cascaderHandleChange" />
                </el-form-item>
                <el-form-item label="Label" class="label-box">
                    <el-tag v-for="tag in form.label" :key="tag" closable :disable-transitions="false" class="label-item"
                        @close="labelHandl.handleClose(tag)">
                        {{ tag }}
                    </el-tag>
                    <el-input v-if="form.labelInputVisible" ref="labelInputRef" v-model="form.labelText" size="small"
                        class="label-input" @keyup.enter="labelHandl.handleInputConfirm"
                        @blur="labelHandl.handleInputConfirm" />
                    <el-button class="label-btn" v-else size="small" @click="labelHandl.showInput">
                        + New Tag
                    </el-button>
                </el-form-item>
                <el-button size="small" type="primary"
                    @click="useSaveData(form, ruleFormRef, router, uploadFileformation, uploadCoveration,props)">{{ props.type ==
                        "edit" ? "Update column" : "Post a column" }}</el-button>
            </el-form>
        </div>
        <!-- 上传组件 -->
        <div class="upload-box" v-show="false">
            <el-upload class="upload" drag action="https://run.mocky.io/v3/9d059bf9-4660-45f2-925d-ce80ad6c4d15"
                ref="uploadRef" multiple :on-success="handle.handleFileSuccess" :on-error="handle.handleFileError"
                :show-file-list="false" :before-upload="handle.beforeFileUpload" :auto-upload="true"
                :http-request="handle.RedefineUploadFile">
                <el-button size="small" type="primary" ref="uploadBtnRef">upload files</el-button>
            </el-upload>
        </div>
        <!--File upload progress -->
        <el-progress ref="uploadProgressRef" v-show="false" :text-inside="true" :stroke-width="16" type="dashboard"
            :percentage="uploadFileformation.progress" :color="colors" />
       <!--Rich text component -->
        <div class="quill-box" v-show="!form.isShow">
            <el-page-header @back="form.isShow = true" class="page-header">
                <template #content>
                    <span class=""> {{ form.title }} </span>
                </template>
            </el-page-header>
            <QuillEditor id="editorId" ref="myQuillEditor" v-model:content="form.content" contentType="html"
                :options="options" theme="snow" />
        </div>
    </div>
</template>

<script setup lang="ts">
import '@vueup/vue-quill/dist/vue-quill.snow.css';
import 'highlight.js/styles/agate.css'
import { QuillEditor } from '@vueup/vue-quill'
import { Plus, Edit } from '@element-plus/icons-vue'
import { vRemoveFocus } from "@/utils/customInstruction/focus"
import { useArticleContributionProp, useHandleFileMethod, useHandleCoverMethod, useInit, userLabelHandlMethod, useSaveData, useRules } from '@/logic/creation/contribute/contributePage/articleContribution';
import { onMounted } from 'vue';


const props = defineProps({
    type: {
        type: String
    }
})

components: {
    QuillEditor
}

const { myQuillEditor, form, options, uploadFileformation, uploadCoveration, uploadRef, uploadBtnRef, uploadProgressRef, colors, labelInputRef, ruleFormRef, router, cascader, cascaderValue, editArticleStore } = useArticleContributionProp()
const handle = useHandleFileMethod(uploadFileformation, form, myQuillEditor, uploadProgressRef)
const handleCover = useHandleCoverMethod(uploadCoveration)
const labelHandl = userLabelHandlMethod(form, labelInputRef)
const rules = useRules()
onMounted(() => {
    useInit(uploadFileformation, uploadCoveration, cascader, form, editArticleStore, props, myQuillEditor, cascaderValue)
})

//Category selection
const cascaderHandleChange = (value: any) => {
    form.classification_id = value[value.length - 1]
    console.log(form.classification_id)
}


</script>

<style scoped lang="scss">
@import "@/assets/styles/views/creation/contribute/contributePage/articleContribution.scss"
</style>
