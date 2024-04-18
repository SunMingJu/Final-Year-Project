<template>
  <div class="overall">
    <pageHeader title="Edit information" icon-nmae="userData"></pageHeader>
    <div class="form-box personal-layout animate__animated animate__slideInRight ">
      <el-form :model="form" :rules="userInfoRules" label-width="120px">
        <el-form-item label="Nick name" prop="username">
          <el-input class="input-name" v-model="form.username" clearable />
        </el-form-item>
        <el-form-item label="gender">
          <el-select v-model="form.gender" placeholder="Please select your gender">
            <el-option label="boys" :value="0" />
            <el-option label="girl" :value="1" />
            <el-option label="Schrodinger's cat" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item label="date of birth">
          <el-col :span="11">
            <el-date-picker
              v-model="form.birth_date"
              type="date"
              placeholder="Please select date of birth"
              style="width: 100%"
            />
          </el-col>
        </el-form-item>
        <el-form-item label="Visible to others">
          <el-switch class="switch-btn" v-model="form.is_visible" />
        </el-form-item>

        <el-form-item label="Signature">
          <el-input
            class="input-hobby"
            v-model="form.signature"
            type="textarea"
            resize="none"
            rows="4"
          />
        </el-form-item>
        <el-form-item>
          <el-button v-removeFocus type="primary" round @click="UserInfoMethod.onSubmit">
            Modify information</el-button
          >
        </el-form-item>
      </el-form>
      <div class="figure-box">
        <div class="figure"></div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {
  useUserInfoProp,
  useUserInfoMethod,
  useRules,
} from "@/logic/personal/userInfo/userInfo";
import pageHeader from "@/components/personalNavigation/pageHeader.vue";
import {vRemoveFocus} from "@/utils/customInstruction/focus"
import { onMounted } from "vue";


components: {
  pageHeader;
}
const { form } = useUserInfoProp();
const UserInfoMethod = useUserInfoMethod(form);
const { userInfoRules } = useRules();

onMounted(() => {
  UserInfoMethod.getUserInfo();

})


</script>
<style scoped lang="scss">
@import "@/assets/styles/views/personal/userInfo/userInfo.scss";
@import "@/assets/styles/global/el-style.scss";
</style>
