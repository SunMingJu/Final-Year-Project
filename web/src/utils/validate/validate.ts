import {determineNameExistsRequist} from "@/apis/personal"
import {DetermineNameExistsReq} from "@/types/personal/userInfo/userInfo"
//Verify verification code
export const validateUsername = (rule: any, value: any, callback: any) => {
    value as string;
    console.log(value)
    if (value === '') {
      callback(new Error('please enter user name'));
    } else if (value.length < 2 || value.length > 8) {
      callback(new Error('Nickname length needs to be between 2 and 8 characters'));
    } else{
      callback();
    }
  };
  //verify password
  export const validatePassword = (rule: any, value: any, callback: any) => {
    const reg = new RegExp(/[a-zA-Z0-9!?.]+/);
    if (value === '') {
      callback(new Error('Please enter password'));
    } else if (value.length < 6 || value.length > 16) {
      callback(new Error('The password length needs to be between 6 and 16 characters'));
    } else if (!reg.test(value)) {
      callback(new Error('The password cannot contain special symbols except !?.'));
    } else {
      callback();
    }
  }; 
  //Verify verification code
  export const validateVarCode = (rule: any, value: any, callback: any) => {
    value as string;
    if (value === '') {
      callback(new Error('please enter verification code'));
    } else if (value.length != 6) {
      callback(new Error('The verification code is 6 digits'));
    } else {
      callback();
    }
  };
  //Verify email
  export const validateEmail = (rule: any, value: any, callback: any) => {
    const reg = new RegExp(/[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+/);
    if (value.length === 0) {
      callback(new Error('Can you tell me your email address?~'));
    } else if (value.length < 5 || value.length > 64) {
      callback(new Error('There is something wrong with the length of your email address~'));
    } else if (!reg.test(value)) {
      callback(new Error('Can you enter the correct email format?~'));
    } else {
      callback();
    }
  };


   //Verify duplicate nickname
   export const validateRepeatName = async (rule: any, value: any, callback: any) => {
    value as string;
    let requist = <DetermineNameExistsReq>{
      username : value
    }
    console.log(value)
    if (value === '') {
      callback(new Error('please enter user name'));
    } else if (value.length < 2 || value.length > 8) {
      callback(new Error('Nickname length needs to be between 2 and 16 characters'));
    } else if ((await determineNameExistsRequist(requist)).data) {
      callback(new Error('Someone has used the nickname.'));
    } else{
      callback();
    }
};


   //Verify live broadcast title
   export const validateLiveTitle = async (rule: any, value: any, callback: any) => {
    value as string;
    if (value === '') {
      callback(new Error('Please enter the live broadcast title'));
    } else if (value.length < 8 || value.length > 20) {
      callback(new Error('Title length needs to be between 8 and 20 characters'));
    } else{
      callback();
    }
};

  //Verify video title
  export const validateVideoTitle = async (rule: any, value: any, callback: any) => {
    value as string;
    if (value === '') {
      callback(new Error('Please enter the live broadcast title'));
    } else if (value.length < 4 || value.length > 30) {
      callback(new Error('Title length needs to be between 4 and 30 characters'));
    } else{
      callback();
    }
};

  //Verification video introduction
  export const validateVideoIntroduce = async (rule: any, value: any, callback: any) => {
    value as string;
    if (value === '') {
      callback(new Error('Please enter live broadcast introduction'));
    }else{
      callback();
    }
};

  //Verify video title
  export const validateArticleTitle = async (rule: any, value: any, callback: any) => {
    value as string;
    if (value === '') {
      callback(new Error('Please enter the article title'));
    } else if (value.length < 8 || value.length > 20) {
      callback(new Error('Title length needs to be between 8 and 20 characters'));
    } else{
      callback();
    }
};



   //Verify favorites title
   export const validateCollectTitle = async (rule: any, value: any, callback: any) => {
    value as string;
    if (value === '') {
      callback(new Error('Please enter a title'));
    } else if (value.length < 1 || value.length > 20) {
      callback(new Error('Title length needs to be between 1 and 20 characters'));
    } else{
      callback();
    }
};