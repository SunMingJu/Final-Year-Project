<template>
    <div class="overall">
        <div class=" principal  personal-layout">
            <div class="form-box">
                <el-upload class="cover-uploader" :action="createFavoriteRmationForm.action" :show-file-list="false"
                    :on-success="handle.handleFileSuccess" :on-error="handle.handleFileError"
                    :before-upload="handle.beforeFileUpload" :auto-upload="true" :http-request="handle.RedefineUploadFile"
                    accept=".png,.jpg,.jpeg">
                    <img v-if="createFavoriteRmationForm.FileUrl" :src="createFavoriteRmationForm.FileUrl" class="cover" />
                    <el-icon v-else class="cover-uploader-icon">
                        <Plus />
                    </el-icon>
                </el-upload>
                <div>
                    <div class="form-show">
                        <el-form :model="createFavoriteRmationForm" ref="saveDateFormRef" :rules="liveInformationRules"
                            label-position="left" label-width="5rem">
                            <el-form-item label="title" prop="title">
                                <el-input v-model="createFavoriteRmationForm.title" />
                            </el-form-item>
                            <el-form-item label="introduce">
                                <el-input type="textarea" resize="none" :rows="4"
                                    v-model="createFavoriteRmationForm.content" />
                            </el-form-item>
                        </el-form>
                    </div>
                </div>
            </div>
            <div class="bottom-box">
                <span class="text"> Please set your cover and title to better search for collection content.~ </span>
                <div class="button">
                    <el-button v-removeFocus
                        @click="useSaveData(createFavoriteRmationForm, saveDateFormRef, rawData, router, emits)"
                        type="primary" round> {{createFavoriteRmationForm.id == 0 ? "Create folder" : "update folder"}}
                    </el-button>
                </div>
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>

import { Plus } from "@element-plus/icons-vue";
import { vRemoveFocus } from "@/utils/customInstruction/focus"
import { useHandleFileMethod, useSaveData, useRules, useInit, useCreateFavoritesProp } from "@/logic/personal/create/createFavorites"
import { onMounted, onUpdated } from "vue";
import { GetFavoritesListItem } from "@/types/personal/collect/favorites";

const props = defineProps({
    info: {
        type: Object as () => GetFavoritesListItem,
    },
    type : {
        type : Boolean, //true is insert mode false is update mode
        required: true,
    }
})
const emits = defineEmits(["shutDown"])

const { createFavoriteRmationForm, saveDateFormRef, rawData, router } = useCreateFavoritesProp()
const handle = useHandleFileMethod(createFavoriteRmationForm)
const { liveInformationRules } = useRules();

onMounted(() => {
    useInit(createFavoriteRmationForm, rawData,props.type,props.info)

})

//Data changes are updated again
onUpdated(() => {
    useInit(createFavoriteRmationForm, rawData,props.type,props.info)
})
</script>

<style lang="scss" scoped>@import "./static/style/createFavorites.scss";</style>
