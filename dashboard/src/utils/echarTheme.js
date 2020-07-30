


import echarts from 'echarts';

var colorPalette = [
    '#1B2948', '#759aa0', '#e69d87', '#8dc1a9', '#ea7e53',
    '#eedd78', '#73a373', '#73b9bc', '#7289ab', '#91ca8c', '#f49f42'
];

echarts.registerTheme('my_theme', {
    color: colorPalette,
    //backgroundColor: '#030417'
  });