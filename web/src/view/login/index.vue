<template>
  <div id="userLayout" class="w-full h-full min-h-screen relative overflow-hidden">
    <!-- 背景图片层 -->
    <img src="@/assets/login-bg.png" alt="login background" class="absolute inset-0 w-full h-full object-cover z-0">
    <!-- 遮罩层（可选，用于提高内容可读性） -->
    <div class="absolute inset-0 bg-black bg-opacity-10 z-0"></div>
    
    <div class="shadow-2xl rounded-lg flex flex-col md:flex-row w-[52%] h-[57%] absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 overflow-hidden z-10">
      <!-- 左侧品牌区域 -->
      <div class="w-[50%] bg-[#4699F5] flex flex-col items-start p-15 relative overflow-hidden">
        <div class="relative z-10 w-full flex flex-col gap-6 h-full items-start">
          <img class="h-[3.75rem] w-auto object-contain" src="@/assets/logo-ln.png" alt="logo">
          <p class="text-white text-3xl mb-8">四可边缘(AGC&AVC)自动控制系统</p>
          <img src="@/assets/login-2.png" alt="login-2" class="max-h-[45%] w-auto object-contain">
        </div>
      </div>
      
      <!-- 右侧登录区域 -->
      <div class="w-[50%] bg-white flex items-center p-20 box-border">
        <div class="w-full max-w-md flex flex-col gap-12">
          <h2 class="text-[2.75rem] font-thin mb-8">欢迎登录</h2>
          
          <el-form
            ref="loginForm"
            :model="loginFormData"
            :rules="rules"
            :validate-on-rule-change="false"
            @keyup.enter="submitForm"
            class="space-y-8"
            size="medium"
          >
            <el-form-item prop="username">
              <el-input
                v-model="loginFormData.username"
                size="large"
                placeholder="请输入用户名"
                class="border-gray-300 focus:border-blue-500 focus:ring-blue-500"
              />
            </el-form-item>
            
            <el-form-item prop="password">
              <el-input
                v-model="loginFormData.password"
                show-password
                size="large"
                type="password"
                placeholder="请输入密码"
                class="border-gray-300 focus:border-blue-500 focus:ring-blue-500"
              />
            </el-form-item>
            
            <!-- <el-form-item v-if="loginFormData.openCaptcha" prop="captcha">
              <div class="flex w-full gap-4">
                <el-input
                  v-model="loginFormData.captcha"
                  placeholder="请输入验证码"
                  size="large"
                  prefix-icon="el-icon-document-checked"
                  class="flex-1 border-gray-300 focus:border-blue-500 focus:ring-blue-500"
                />
                <div class="w-1/3 h-12 bg-[#f0f2f5] rounded-md overflow-hidden">
                  <img
                    v-if="picPath"
                    class="w-full h-full object-cover cursor-pointer hover:opacity-90 transition-opacity"
                    :src="picPath"
                    alt="请输入验证码"
                    @click="loginVerify()"
                  />
                </div>
              </div>
            </el-form-item> -->
            
            <el-form-item class="flex items-center justify-between text-blue-500">
              <el-checkbox v-model="loginFormData.rememberMe" >记住密码</el-checkbox>
            </el-form-item>
            
            <el-form-item>
              <el-button
                class="w-full mt-8 h-12 text-white bg-blue-600 hover:bg-blue-700 focus:ring-4 focus:ring-blue-300 transition-all duration-300"
                type="primary"
                size="large"
                @click="submitForm"
              >登录</el-button>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </div>
    
    <!-- 页脚信息 -->
    <div class="absolute bottom-0 left-0 right-0 text-center text-white text-sm py-6 z-10">
      <p>Copyright © 2018-2025 绿能中环（西安）新能源科技有限公司 All Rights Reserved</p>
    </div>
  </div>
</template>

<script setup>
  import { captcha } from '@/api/user'
  import { checkDB } from '@/api/initdb'
  import { reactive, ref, onMounted } from 'vue'
  import { ElMessage } from 'element-plus'
  import { useRouter } from 'vue-router'
  import { useUserStore } from '@/pinia/modules/user'

  defineOptions({
    name: 'Login'
  })

  const router = useRouter()
  
  // 验证函数
  const checkUsername = (rule, value, callback) => {
    if (!value || value.length < 5) {
      return callback(new Error('请输入正确的用户名'))
    } else {
      callback()
    }
  }
  
  const checkPassword = (rule, value, callback) => {
    if (!value || value.length < 6) {
      return callback(new Error('请输入正确的密码'))
    } else {
      callback()
    }
  }

  // 获取验证码
  const loginVerify = async () => {
    try {
      const ele = await captcha()
      if (ele.data) {
        rules.captcha = [{
          max: ele.data.captchaLength,
          min: ele.data.captchaLength,
          message: `请输入${ele.data.captchaLength}位验证码`,
          trigger: 'blur'
        }]
        picPath.value = ele.data.picPath
        loginFormData.captchaId = ele.data.captchaId
        loginFormData.openCaptcha = ele.data.openCaptcha
      }
    } catch (error) {
      console.error('获取验证码失败:', error)
      ElMessage({
        type: 'error',
        message: '获取验证码失败，请稍后重试',
        showClose: true
      })
    }
  }

  // 登录相关操作
  const loginForm = ref(null)
  const picPath = ref('')
  const loginFormData = reactive({
    username: 'admin',
    password: '',
    captcha: '',
    captchaId: '',
    openCaptcha: false,
    rememberMe: false
  })
  
  const rules = reactive({
    username: [{ validator: checkUsername, trigger: 'blur' }],
    password: [{ validator: checkPassword, trigger: 'blur' }],
    captcha: []
  })

  const userStore = useUserStore()
  
  const login = async () => {
    try {
      return await userStore.LoginIn(loginFormData)
    } catch (error) {
      console.error('登录失败:', error)
      ElMessage({
        type: 'error',
        message: '登录失败，请检查网络或联系管理员',
        showClose: true
      })
      return false
    }
  }
  
  const submitForm = () => {
    if (!loginForm.value) return
    
    loginForm.value.validate(async (v) => {
      if (!v) {
        // 未通过前端静态验证
        ElMessage({
          type: 'error',
          message: '请正确填写登录信息',
          showClose: true
        })
        // await loginVerify()
        return false
      }

      // 通过验证，请求登陆
      const flag = await login()

      // 登陆失败，刷新验证码
      if (!flag) {
        // await loginVerify()
        return false
      }

      // 登陆成功
      return true
    })
  }

  // 初始化加载
  onMounted(() => {
    // loginVerify()
  })
</script>

<style scoped>
  #userLayout {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  }
  
  /* 自定义滚动条样式 */
  ::-webkit-scrollbar {
    width: 6px;
    height: 6px;
  }
  
  ::-webkit-scrollbar-thumb {
    background-color: rgba(144, 147, 153, 0.3);
    border-radius: 3px;
  }
  
  ::-webkit-scrollbar-track {
    background-color: transparent;
  }
</style>
