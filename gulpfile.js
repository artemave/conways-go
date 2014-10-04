var gulp           = require('gulp')
var pogo           = require('gulp-pogo')
var browserify     = require('gulp-browserify');
var sass           = require('gulp-sass')
var concat         = require('gulp-concat');
var plumber        = require('gulp-plumber')
var gutil          = require('gulp-util')
var fs             = require('fs')
var gulpBowerFiles = require('main-bower-files')
var watch          = require('gulp-watch')
var karma          = require('karma').server;

var onError = function (err) {
  gutil.beep();
  gutil.log(gutil.colors.red(err.message))
  gutil.log(err)
};

gulp.task('styles', function (callback) {
  return gulp.src('./public/css/app.scss')
    .pipe(plumber({
      errorHandler: onError
    }))
    .pipe(sass())
    .pipe(concat('bundle.css'))
    .pipe(gulp.dest('./public'))
});

gulp.task('compile-pogo', function(callback){
    return gulp.src('./public/js/**/*.pogo')
      .pipe(plumber({errorHandler: onError}))
      .pipe(pogo())
      .pipe(gulp.dest('./public/js/'));
})

gulp.task('scripts', ['compile-pogo'], function() {
    return gulp.start('browserify')
});

gulp.task('browserify', function() {
  return gulp.src('./public/js/app.js', {read: false})
    .pipe(plumber({
      errorHandler: onError
    }))
    .pipe(browserify({
      insertGlobals: true
    }))
    .pipe(concat('bundle.js'))
    .pipe(gulp.dest('./public'))
})


gulp.task("bower-files", function() {
  return gulp.src(gulpBowerFiles())
    .pipe(concat('deps.js'))
    .pipe(gulp.dest("./public"))
});

/**
 * Run test once and exit
 */
gulp.task('test', function (done) {
  karma.start({
    configFile: __dirname + '/karma.conf.js',
    singleRun: true
  }, done);
});

/**
 * Watch for file changes and re-run tests on each change
 */
gulp.task('tdd', function (done) {
  karma.start({
    configFile: __dirname + '/karma.conf.js'
  }, done);
});

gulp.task("watch", function() {
  watch('./public/test/**/*.pogo')
    .pipe(plumber({errorHandler: onError}))
    .pipe(pogo())
    .pipe(browserify())
    .pipe(gulp.dest('./public/test/'));

  watch('./public/js/**/*.pogo')
    .pipe(plumber({errorHandler: onError}))
    .pipe(pogo())
    .pipe(gulp.dest('./public/js/'));

  gulp.watch('./public/js/**/*.js', ['browserify']);
  gulp.watch('./public/css/**', ['styles']);
})

gulp.task('default', ['styles', 'bower-files', 'scripts']);
