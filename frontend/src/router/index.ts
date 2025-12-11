// @ts-nocheck
import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';
import {useAppStore} from "@/store/app";
import Panel from '../pages/User/Panel.vue';
import Login from '../pages/User/Login.vue';
import Home from '../pages/User/Home.vue';
import Poems from "../pages/User/Poems.vue";
import LikedPoems from "../pages/User/LikedPoems.vue";
import Register from "../pages/register.vue";
import RegisterUser from "../pages/RegisterUser.vue";
import LatestPoems from "../pages/User/LatestPoems.vue";
import PopularPoems from "../pages/User/PopularPoems.vue";
import Poem from "../pages/User/Poem.vue";
import Book from "../pages/User/Book/BookList.vue";
import BookDetail from "../pages/User/Book/BookDetail.vue";
import Bookmark from "../pages/User/Bookmark.vue";
import Reminder from "../pages/User/Reminder.vue";
import ReadsBooks from "../pages/User/ReadsBooks.vue";
import BookToReads from "../pages/User/NotReadsBooks.vue";
import Friends from "../pages/User/Friends.vue";
import Authors from "../pages/User/Authors.vue";
import AuthorDetail from "../pages/User/AuthorDetail.vue";
import UserProfile from "../pages/User/UserProfile.vue";
import AdminManagement from '../pages/Management/AdminManagement.vue';
import PoemManagement from "../pages/Management/PoemManagement.vue";
import LogManagement from "../pages/Management/LogManagement.vue";
import BookManagement from "../pages/Management/BookManagement.vue";
import AuthorManagement from '../pages/Management/AuthorManagement.vue';
import ReminderManagement from '../pages/Management/reminderManagement.vue';
import HomepageManagement from '../pages/Management/HomepageManagement.vue';
import MihrimahCardManagement from '../pages/Management/MihrimahCardManagement.vue';
import NotFound from '../pages/NotFound.vue';
import axios from "axios";

const routes: Array<RouteRecordRaw> = [
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/register',
    name: 'RegisterUser',
    component: RegisterUser,
    meta: { requiresAuth: false },
  },
  {
    path: '/panel',
    component: () => import('../layouts/adminPanel.vue'),
    children: [
      {
        path: '/panel',
        name: 'Panel',
        component:Panel,
        meta: { requiresAuth: true },
      },
      {
        path: '/admin-management',
        name: 'AdminManagement',
        component:AdminManagement,
        meta: { requiresAuth: true },
      },
      {
        path: '/poem-management',
        name: 'PoemManagement',
        component:PoemManagement,
        meta: { requiresAuth: true },
      },
      {
        path: '/log-management',
        name: 'LogManagement',
        component:LogManagement,
        meta: { requiresAuth: true },
      },
      {
        path: '/book-management',
        name: 'BookManagement',
        component:BookManagement,
        meta: { requiresAuth: true },
      },
      {
        path: '/author-management',
        name: 'AuthorManagement',
        component: AuthorManagement,
        meta: { requiresAuth: true },
      },
      {
        path: '/reminder-management',
        name: 'ReminderManagement',
        component: ReminderManagement,
        meta: { requiresAuth: true },
      },
      {
        path: '/homepage-management',
        name: 'HomepageManagement',
        component: HomepageManagement,
        meta: { requiresAuth: true },
      },
      {
        path: '/mihrimah-card-management',
        name: 'MihrimahCardManagement',
        component: MihrimahCardManagement,
        meta: { requiresAuth: true },
      },
    ]
  },
  {
    path: '/',
    component: () => import('../layouts/default.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '/home',
        name: 'Home',
        component:Home,
        meta: { requiresAuth: true },
      },
      {
        path: '/poems',
        name: 'Poems',
        component:Poems,
        meta: { requiresAuth: true },
      },
      {
        path: '/latest-poems',
        name: 'LatestPoems',
        component:LatestPoems,
        meta: { requiresAuth: true },
      },
      {
        path: '/popular-poems',
        name: 'PopularPoems',
        component:PopularPoems,
        meta: { requiresAuth: true },
      },
      {
        path: '/poem/:slug',
        name: 'Poem',
        component:Poem,
        meta: { requiresAuth: true },

      },
      {
        path: '/liked-poems',
        name: 'LikedPoems',
        component:LikedPoems,
        meta: { requiresAuth: true },
      },
      {
        path: '/bookmark',
        name: 'Bookmark',
        component:Bookmark,
        meta: { requiresAuth: true },
      },
      {
        path: '/books',
        name: 'Book',
        component:Book,
        meta: { requiresAuth: true },
      },
      {
        path: '/book/:slug',
        name: 'BookDetail',
        component: BookDetail,
        meta: { requiresAuth: true },
        props: true
      },
      {
        path: '/reads-books',
        name: 'ReadsBooks',
        component:ReadsBooks,
        meta: { requiresAuth: true },
      },
      {
        path: '/book-to-reads',
        name: 'BookToReads',
        component:BookToReads,
        meta: { requiresAuth: true },
      },
      {
        path: '/sozler',
        name: 'Sozler',
        component:Reminder,
        meta: { requiresAuth: true },
      },
      {
        path: '/friends',
        name: 'Friends',
        component:Friends,
        meta: { requiresAuth: true },
      },
      {
        path: '/authors',
        name: 'Authors',
        component:Authors,
        meta: { requiresAuth: true },
      },
      {
        path: '/author/:slug',
        name: 'AuthorDetail',
        component:AuthorDetail,
        meta: { requiresAuth: true },
        props: true
      },
      {
        path: '/profile/:username',
        name: 'UserProfile',
        component:UserProfile,
        meta: { requiresAuth: true },
        props: true
      },
      {
        path: ':pathMatch(.*)*',
        name: 'NotFound',
        component: NotFound,
        meta: { requiresAuth: true },
      },
    ]
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

router.beforeEach(async (to, from, next) => {
  const store = useAppStore();
  const forbiddenRoutes = ['Panel', 'AdminManagement', 'PoemManagement', 'HomepageManagement', 'MihrimahCardManagement'];
  const adminMemberOnlyRoutes = ['']; // role_id 1 veya 2 olanlar girebilir

  if (to.path === '/') {
    return next('/home');
  }

  if (to.meta.requiresAuth) {
    try {
      const resAuth = await axios.get('/auth-check', { withCredentials: true });
      if( resAuth.data.message !== 'ok' ) {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        return next('/login');
      }

      // AuthCheck'ten minimal user bilgisi al (id ve role_id)
      const user = resAuth.data.user;
      const role_id = user.role_id;

      // Store'a minimal user bilgisini kaydet
      store.setUser(user);

      if (role_id) {
        // Sadece admin'lerin girebileceği sayfalar (role_id = 1)
        if(forbiddenRoutes.includes(to.name as string) && role_id !== 1) {
          next(`/home`);
          scrollTo(0,0);
        }
        // Sadece admin ve member'ların girebileceği sayfalar (role_id = 1 veya 2)
        else if(adminMemberOnlyRoutes.includes(to.name as string) && role_id !== 1 && role_id !== 2) {
          next(`/home`);
          scrollTo(0,0);
        }
        else{
          next();
          scrollTo(0,0);
        }

      } else {
        // console.log("panelden çıkış yapıldı")
        next('/login');
        scrollTo(0,0);
      }
    } catch (error) {
      // console.log("panelden çıkış yapıldı")
      next('/login');
      scrollTo(0,0);
    }
  } else {
    next();
    scrollTo(0,0);
  }
});

export default router;
