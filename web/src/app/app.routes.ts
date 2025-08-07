import {Routes} from '@angular/router';

export const routes: Routes = [
  {
    title: 'Home',
    path: '',
    loadComponent: () => import('../pages/home/home').then(m => m.Home),
  },
  {
    title: 'About',
    path: 'about',
    loadComponent: () => import('../pages/about/about').then(m => m.About),
  },
  {
    title: 'AdBlock',
    path: 'adblock',
    loadComponent: () => import('../pages/adblock/adblock').then(m => m.Adblock),
  },
  {
    path: '**',
    loadComponent: () => import('../pages/errors/errors').then(m => m.Errors),
  },
];
