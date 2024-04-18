<template>
  <div class="main ">
<!--login/register main box -->
    <div :class="{ active: currentModel, container: true,animate__animated :true, animate__flipInX :true}" >
      <!--Register -->
      <div class="form-container sign-up-container ">
        <el-form
          v-if="!isForget"
          class="form"
          :model="regForm"
          onsubmit="return false;"
          ref="regFormRef"
          :rules="registerRules"
        >
          <h3 class="container-title">Register</h3>
          <el-form-item class="form-item" prop="username">
            <el-input v-model="regForm.username" type="text" placeholder="Nick name" />
          </el-form-item>
          <el-form-item class="form-item" prop="email">
            <el-input
              v-model="regForm.email"
              placeholder="Mail"
              class="input-with-select"
            >
              <template #append v-if="regForm.email">
                <el-button
                  type="primary"
                  size="small"
                  round
                  @click="sendEmail(regForm.email, 'regist')"
                  >{{ sendEmailInfo.btnText }}</el-button
                >
              </template>
            </el-input>
          </el-form-item>
          <el-form-item class="form-item" prop="password">
            <el-input
              v-model="regForm.password"
              type="password"
              name="password"
              placeholder="password"
            />
          </el-form-item>
          <el-form-item class="form-item" prop="verificationCode">
            <el-input
              v-model="regForm.verificationCode"
              type="text"
              name="password"
              placeholder="Verification code"
            />
          </el-form-item>
          <button  type="button" class="signUp theme-button" @click="register(regFormRef)">sign up</button>
        </el-form>

        <!--Retrieve password -->
        <el-form
          v-if="isForget"
          class="form"
          :model="forgetForm"
          onsubmit="return false;"
          ref="forgetFormRef"
          :rules="forgetRules"
        >
          <h3 class="container-title">Forget</h3>
          <el-form-item class="form-item" prop="email">
            <el-input
              v-model="forgetForm.email"
              placeholder="Mail"
              class="input-with-select"
            >
              <template #append v-if="forgetForm.email">
                <el-button
                  type="primary"
                  size="small"
                  round
                  @click="sendEmail(forgetForm.email, 'modify')"
                  >{{ sendEmailInfo.btnText }}</el-button
                >
              </template>
            </el-input>
          </el-form-item>
          <el-form-item class="form-item" prop="password">
            <el-input
              v-model="forgetForm.password"
              type="password"
              name="password"
              placeholder="password"
            />
          </el-form-item>
          <el-form-item class="form-item" prop="verificationCode">
            <el-input
              v-model="forgetForm.verificationCode"
              type="text"
              name="password"
              placeholder="Verification code"
            />
          </el-form-item>
          <button   type="button" class="signUp theme-button" @click="forfet(forgetFormRef)">modify</button>
        </el-form>
      </div>
      <!--Login -->
      <div class="form-container sign-in-container">
        <el-form
          class="form"
          :model="loginForm"
          onsubmit="return false;"
          :rules="loginRules"
          ref="loginFormRef"
        >
          <h2 class="container-title">Login</h2>

          <el-form-item class="form-item" prop="username">
            <el-input v-model="loginForm.username" type="text" placeholder="Nick name" />
          </el-form-item>
          <el-form-item class="form-item" prop="passwoed">
            <el-input
              v-model="loginForm.password"
              type="password"
              name="password"
              placeholder="password"
            />
          </el-form-item>
          <button class="signUp theme-button"   type="button"  @click="login(loginFormRef)">sign in</button>
          <p
            class="forget"
            @click="
              currentModel = !currentModel;
              isForget = true;
            "
          >
            --- Forget Passwoed----
          </p>
        </el-form>
      </div>
      <!--Overlay container-->
      <div class="overlay_container">
        <div class="overlay">
          <!--Left side -->
          <div class="overlay_panel overlay_left_container">
            <h2 class="container-title">Welcome back!</h2>
            <p>Please log in with your personal information to stay in touch with us</p>
            <button  type="button" id="sign-in" class="theme-button" @click="currentModel = !currentModel">login</button>
          </div>
          <!--Right -->
          <div class="overlay_panel overlay_right_container">
            <h2 class="container-title">hello friend!</h2>
            <p>Enter your personal information and start the journey with us</p>
            <button
            type="button"
            class="theme-button"
            id="sign-up"
              @click="
                currentModel = !currentModel;
                isForget = false;
              "
            >
              register
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {
  useLoginProp,
  useLoginMethod,
  useForgetrMethod,
  useRules,
  useRegisterMethod,
} from "@/logic/login/login";

const {
  currentModel,
  userStore,
  router,
  regForm,
  regFormRef,
  loginForm,
  loginFormRef,
  forgetForm,
  sendEmailInfo,
  sendEmail,
  isForget,
  forgetFormRef,
} = useLoginProp();

console.log(currentModel);

const { loginRules, registerRules, forgetRules } = useRules();
const { login } = useLoginMethod(userStore, router, loginForm);
const { register } = useRegisterMethod(userStore, router, regForm);
const { forfet } = useForgetrMethod(forgetForm, currentModel);
</script>

<style scoped lang="scss">
@import "@/assets/styles/views/login/login.scss";
</style>
