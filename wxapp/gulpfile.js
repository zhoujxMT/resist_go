var gulp = require('gulp')
var del = require('del')
var path = require('path')
var ts = require('gulp-typescript')
var jsonminify = require('gulp-jsonminify2')
var htmlmin = require('gulp-htmlmin')
var sass = require('gulp-sass-china')
var uglify = require('gulp-uglify')
const autoprefixer = require('gulp-autoprefixer')
var tsProject = ts.createProject("tsconfig.json")
var minifycss = require('gulp-minify-css')

// ------------ 测试任务集 --------------
gulp.task("json", ()=>{
    return gulp.src('./src/**/*.json')
    .pipe(gulp.dest('./dist'))
})


gulp.task("wxml", ()=>{
    return gulp.src('./src/**/*.wxml').
    pipe(gulp.dest('./dist'))
})

gulp.task("wxss", ()=>{
    return gulp.src('./src/**/*.{wxss, sass}').
    sass().on('error', sass.logError).
    autoprefixer([
        'iOS >=8',
        'Android >=4.1'
    ]).
    rename((path)=>path.extname='.wxss').
    pipe(gulp.dest('./dist'))
})

gulp.task("js", ()=>{
    return gulp.src('./src/**/*.js').
    pipe(gulp.dest('./dist'))
})

gulp.task("ts", ()=> {
    return gulp.src('./src/**/*.ts').
    pipe(tsProject()).js.
    pipe(gulp.dest("./dist"))
})
// ------------------------------------------
// ------------ 发布任务集 --------------
gulp.task("jsonPro", ()=>{
    return gulp.src('./src/**/*.json').
    pipe(jsonminify())
    .pipe(gulp.dest('./dist'))
})


gulp.task("wxmlPro", ()=>{
    return gulp.src('./src/**/*.wxml').
    pipe(htmlmin({
        collapseWhitespace:true,
        removeComments: true,
        keepClosingSlash:true
    })).
    pipe(gulp.dest('./dist'))
})

gulp.task("wxssPro", ()=>{
    return gulp.src('./src/**/*.{wxss, sass}').
    sass().on('error', sass.logError).
    autoprefixer([
        'iOS >=8',
        'Android >=4.1'
    ]).
    minifycss().
    rename((path)=> path.extname='.wxss').
    pipe(gulp.dest('./dist'))
})

gulp.task("jsPro", ()=>{
    return gulp.src('./src/**/*.js').
    pipe(gulp.dest('./dist'))
})

gulp.task("tsPro", ()=> {
    return gulp.src('./src/**/*.ts').
    pipe(tsProject()).js.
    pipe(uglify({
        compress:true
    }))
    pipe(gulp.dest("./dist"))
})

//命令
gulp.task("dev",["json","wxml", "wxss","js", "ts"])
gulp.task("clean", ()=> {
    return del(['./dist**'])
})
gulp.task("build", ["jsonPro", "wxmlPro","wxssPro", "jsPro", "tsPro"])
gulp.task("default", ["dev"])
gulp.task("watch",()=>{
    gulp.watch('./src/**/*.json',["json"]);
    gulp.watch('./src/**/*.ts',['ts']);
    gulp.watch('./src/**/*.js',['js']);
    gulp.watch('./src/**/*.wxml',['wxml']);
    gulp.watch('./src/**/*.sass',['wxss']);
    gulp.watch('./src/**/*.wxss', ['wxss']);
})
