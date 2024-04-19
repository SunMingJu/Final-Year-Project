import { forgetRequist, loginRequist, regist, sendEmailVerificationCode, sendEmailVerificationCodeByForget } from "@/apis/login";
import globalScss from "@/assets/styles/global/export.module.scss";
import { useGlobalStore, useUserStore } from "@/store/main";
import { Result } from "@/types/idnex";
import { forgetReq, loginReq, registReq, sendEmailInfo, sendEmailReq, sendEmailType, userInfoRes } from "@/types/login/login";
import { validateEmail, validatePassword, validateUsername, validateVarCode } from "@/utils/validate/validate";
import { FormInstance } from "element-plus";
import Swal from "sweetalert2";
import { Ref, reactive, ref } from "vue";
import { Router, useRouter } from 'vue-router';
import { useInitChatSocket } from "../global/chat";
import { useInitNoticeSocket } from "../global/notice";
//provide data

const loading = useGlobalStore().globalData.loading

const Toast = Swal.mixin({
    toast: true,
    position: 'top',
    showConfirmButton: false,
    timer: 3000,
})

export const useLoginProp = () => {
    //Public section
    const currentModel = ref(false) //Current module true registration false login
    const userStore = useUserStore()
    const router = useRouter();
    const isShowSendEmail = ref(false)
    const loginFormRef = ref<FormInstance>()
    const isForget = ref(false)


    const loginForm = reactive<loginReq>({
        username: "",
        password: ""
    })

    const regFormRef = ref<FormInstance>()
    const regForm = reactive<registReq>({
        email: "",
        username: "",
        password: "",
        verificationCode: "",
    })

    const forgetFormRef = ref<FormInstance>()
    const forgetForm = reactive<forgetReq>({
        email: "",
        password: "",
        verificationCode: "",
    })

    let sendEmailInfo = reactive<sendEmailInfo>({
        btnText: "send",
        isPleaseClick: true
    })

    const sendEmail = async (email: string, emailType: sendEmailType) => {
        if (!emailType) return
        if (sendEmailInfo.isPleaseClick == false) return
        try {
            const emailReq = <sendEmailReq>{
                email: email
            }
            loading.loading = true
            if (emailType == "regist") {
                await sendEmailVerificationCode(emailReq)
            } else {
                await sendEmailVerificationCodeByForget(emailReq)
            }

            sendEmailInfo.isPleaseClick = false
            loading.loading = false
            Toast.fire({
                icon: 'success',
                title: 'Verification code sent successfully'
            })
            sendEmailInfo.btnText = "60"
            const interval = setInterval(() => {
                if (Number(sendEmailInfo.btnText) <= 0) {
                    sendEmailInfo.btnText = "send"
                    sendEmailInfo.isPleaseClick = true
                    clearInterval(interval)
                } else {
                    sendEmailInfo.btnText = String(Number(sendEmailInfo.btnText) - 1)
                }
            }, 1000)


        } catch (err) {
            loading.loading = false
            console.log(err)

        }
    }

    return {
        currentModel,
        userStore,
        router,
        loginForm,
        forgetForm,
        loginFormRef,
        regForm,
        regFormRef,
        isShowSendEmail,
        sendEmailInfo,
        sendEmail,
        isForget,
        forgetFormRef
    }
}
//login
export const useLoginMethod = (store: any, router: Router, loginForm: loginReq) => {

    const login = async (formEl: FormInstance | undefined) => {
        console.log(formEl)
        if (!formEl) return;
        await formEl.validate(async (valid, fields) => {
            if (valid) {
                try {
                    loading.loading = true
                    const result = await loginRequist(loginForm)
                    store.setUserInfo(<userInfoRes>result.data)
                    //Initialize socket
                    useInitChatSocket()
                    useInitNoticeSocket()
                
                    loading.loading = false
                    router.push("/")

                }
                catch (err: any) {
                    err as unknown as Result
                    loading.loading = false
                    Swal.fire({
                        title: err.message,
                        heightAuto: false,
                        confirmButtonColor: globalScss.colorButtonTheme,
                        icon: "error",
                    })
                    console.log(err)
                }
            } else {
                console.log('error submit!', fields)
            }
        })

    }
    return {
        login
    }
}


//register
export const useRegisterMethod = (store: any, router: Router, registForm: registReq) => {

    const register = async (formEl: FormInstance | undefined) => {
        console.log(formEl)
        if (!formEl) return;
        await formEl.validate(async (valid, fields) => {
            if (valid) {
                try {
                    loading.loading = true
                    const result = await regist(registForm)
                    store.setUserInfo(<userInfoRes>result.data)
                    Toast.fire({
                        icon: 'success',
                        title: '注册成功'
                    })
                    loading.loading = false
                    router.push("/")
                }
                catch (err: any) {
                    err as unknown as Result
                    loading.loading = false
                    Swal.fire({
                        title: err.message,
                        heightAuto: false,
                        confirmButtonColor: globalScss.colorButtonTheme,
                        icon: "error",
                    })
                    console.log(err)
                }
            } else {
                console.log('error submit!', fields)
            }
        })
    }
    return {
        register
    }
}


//Retrieve password
export const useForgetrMethod = (forgetForm: forgetReq, currentModel: Ref<boolean>) => {

    const forfet = async (formEl: FormInstance | undefined) => {
        if (!formEl) return;
        await formEl.validate(async (valid, fields) => {
            if (valid) {
                try {
                    loading.loading = true
                    await forgetRequist(forgetForm)
                    Toast.fire({
                        icon: 'success',
                        title: 'Successfully modified'
                    })
                    forgetForm.email = ""
                    forgetForm.password = ""
                    forgetForm.verificationCode = ""
                    //Clear form validation requires delay
                    setTimeout(() => {
                        formEl.clearValidate()
                    }, 500)
                    currentModel.value = false
                    loading.loading = false
                    return

                }
                catch (err: any) {
                    err as unknown as Result
                    loading.loading = false
                    Swal.fire({
                        title: err.message,
                        heightAuto: false,
                        confirmButtonColor: globalScss.colorButtonTheme,
                        icon: "error",
                    })
                    console.log(err)
                }
            } else {
                formEl.clearValidate()
                console.log('error submit!', fields)
            }
        })
    }
    return {
        forfet
    }
}

//Send the verification code
export const useSendEmail = async (email: string) => {
    try {
        const emailReq = <sendEmailReq>{
            email: email
        }
        const result = await sendEmailVerificationCode(emailReq)
        return result
    } catch (err) {
        console.log(err)
        return err

    }
}

//form validation
export const useRules = () => {
    const loginRules = reactive({
        username: [{ validator: validateUsername, trigger: 'change' }],
        password: [{ validator: validatePassword, trigger: 'change' }]
    });
    // Registration form validation rules
    const registerRules = reactive({
        username: [{ validator: validateUsername, trigger: 'change' }],
        email: [{ validator: validateEmail, trigger: 'change' }],
        password: [{ validator: validatePassword, trigger: 'change' }],
        verificationCode: [{ validator: validateVarCode, trigger: 'change' }]
    });

    const forgetRules = reactive({
        email: [{ validator: validateEmail, trigger: 'change' }],
        password: [{ validator: validatePassword, trigger: 'change' }],
        verificationCode: [{ validator: validateVarCode, trigger: 'change' }]
    });
    return {
        loginRules,
        registerRules,
        forgetRules
    };
};
