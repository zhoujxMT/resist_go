var gulp = require('gulp')
var del = require('del')
var path = require('path')
var ts = require('gulp-typescript')
var jsonminify = require('gulp-jsonminify2')
var htmlmin = require('gulp-htmlmin')
var sass = require('gulp-sass')
var uglify = require('gulp-uglify')
const autoprefixer = require('gulp-autoprefixer')
var tsProject = ts.createProject("tsconfig.json")
var minifycss = require('gulp-minify-css')
const combiner = require('stream-combiner2');
const rename = require('gulp-rename')
const imageMin = require('gulp-imagemin')

const handleError = function(err) {
  console.log('\n')
  gutil.log(colors.red('Error!'))
  gutil.log('fileName: ' + colors.red(err.fileName))
  gutil.log('lineNumber: ' + colors.red(err.lineNumber))
  gutil.log('message: ' + err.message)
  gutil.log('plugin: ' + colors.yellow(err.plugin))
};
// ------------ 测试任务集 --------------
gulp.task("json", () => {
    return gulp.src('./src/**/*.json')
        .pipe(gulp.dest('./dist'))
})


gulp.task("wxml", () => {
    return gulp.src('./src/**/*.wxml').
        pipe(gulp.dest('./dist'))
});

gulp.task("assets", ()=>{
    return gulp.src('./src/assets/**/*.*')
    .pipe(gulp.dest('./dist/assets'))
})
gulp.task("image", ()=>{
    return gulp.src('./src/images/*.png').
    pipe(gulp.dest('./dist/images'))
})

gulp.task("wxss", () => {
    var combined = combiner.obj([
        gulp.src(['./src/**/*.{wxss,scss}', '!./src/styles/**']),
        sass().on('error', sass.logError),
        autoprefixer([
            'iOS >= 8',
            'Android >= 4.1'
        ]),
        rename((path) => path.extname = '.wxss'),
        gulp.dest('./dist')
    ]);
    combined.on('error', handleError);
})

gulp.task("js", () => {
    return gulp.src('./src/**/*.js').
        pipe(gulp.dest('./dist'))
})

gulp.task("ts", () => {
    return gulp.src('./src/**/*.ts').
        pipe(tsProject()).js.
        pipe(gulp.dest("./dist"))
})
// ------------------------------------------
// ------------ 发布任务集 --------------
gulp.task("jsonPro", () => {
    return gulp.src('./src/**/*.json').
        pipe(jsonminify())
        .pipe(gulp.dest('./dist'))
})


gulp.task("wxmlPro", () => {

    return gulp.src('./src/**/*.wxml').
        pipe(htmlmin({
            collapseWhitespace: true,
            removeComments: true,
            keepClosingSlash: true
        })).
        pipe(gulp.dest('./dist'))
})

gulp.task("wxssPro", () => {
    var combined = combiner.obj([
        gulp.src(['./src/**/*.{wxss,scss}', '!./src/styles/**']),
        sass().on('error', sass.logError),
        autoprefixer([
            'iOS >= 8',
            'Android >= 4.1'
        ]),
        rename((path) => path.extname = '.wxss'),
        gulp.dest('./dist')
    ]);
    combined.on('error', handleError);
})

gulp.task("jsPro", () => {
    return gulp.src('./src/**/*.js').
        pipe(gulp.dest('./dist'))
})

gulp.task("tsPro", () => {
    return gulp.src('./src/**/*.ts').
        pipe(tsProject()).js.
        pipe(uglify({
            compress: true
        }))
    pipe(gulp.dest("./dist"))
})

gulp.task("imageMin", ()=>{
    return gulp.src('./src/images/*.*').
    pipe(imageMin({progressive: true}))
    pipe(gulp.dest('./dist/images'))
})

//命令
gulp.task("dev", ["json", "wxml", "wxss", "js", "ts","assets","image"])
gulp.task("clean", () => {
    return del(['./dist**'])
})
gulp.task("build", ["jsonPro", "wxmlPro", "wxssPro", "jsPro", "tsPro","assets","imageMin"])
gulp.task("default", ["dev"])
gulp.task("watch", () => {
    gulp.watch('./src/**/*.json', ["json"]);
    gulp.watch('./src/**/*.ts', ['ts']);
    gulp.watch('./src/**/*.js', ['js']);
    gulp.watch('./src/**/*.wxml', ['wxml']);
    gulp.watch('./src/**/*.scss', ['wxss']);
    gulp.watch('./src/**/*.wxss', ['wxss']);
    gulp.watch('./src/images/*.*', ['image']);
    gulp.watch('./src/assets/**/*.*',['assets'])
})
