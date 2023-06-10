import { createRouter, createWebHashHistory } from 'vue-router';
import LandfillPage from './LandfillPage.vue';
import SurveyPage from './SurveyPage.vue';
import SignInPage from './SignInPage.vue';
import RegionPage from './RegionPage.vue';

const routes = [
  {
    path: '/sign-in',
    name: 'sign-in',
    component: SignInPage,
  },
  {
    path: '/landfills/:id',
    name: 'landfill',
    component: LandfillPage,
  },
  {
    path: '/surveys/:id',
    name: 'survey',
    component: SurveyPage,
  },
  {
    path: '/regions/:id',
    name: 'region',
    component: RegionPage,
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: SignInPage,
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

export default router;
