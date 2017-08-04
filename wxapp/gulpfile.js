var gulp = require('gulp')
var ts = require('gulp-typescript')
var tsProject = ts.createProject("tsconfig.json")

gulp.task("json", ()=>{
    return gulp.src('./src/**/*.json')
    .pipe(gulp.dest('./dist'))
})

gulp.task("wxml", ()=>{
    return gulp.src('./src/**/*.wxml').
    pipe(gulp.dest('./dist'))
})

gulp.task("wxss", ()=>{
    return gulp.src('./src/**/*.wxss').
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

gulp.task("build",["json","wxml", "wxss","js", "ts"])

