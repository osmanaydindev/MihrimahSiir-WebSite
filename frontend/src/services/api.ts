import axios from 'axios';

// Legacy function - keeping for backward compatibility
export const checkAdminInfo = async(token) => {
  // console.log("checkAdminInfo: ", token)
  // console.log(typeof token)
  try{
    const res = await axios.post('/user', {
      data: token
    });
    // console.log("user data: ", res.data.admin)


    // if(!localStorage.getItem('user') ){
    //   localStorage.setItem('user', JSON.stringify(res.data.admin))
    // }
    // console.log("localstorage zaten varmis")
  }catch (e) {
    // console.log(e)
  }
}

// Export all services from the new service layer
export { apiClient } from './api/client'
export { poemService } from './api/poemService'
export { userService } from './api/userService'
export { friendService } from './api/friendService'
export { bookService } from './api/bookService'
export { commentService } from './api/commentService'
export { authorService } from './api/authorService'
