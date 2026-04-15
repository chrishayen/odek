# Requirement: "a low-dimensional linear algebra library for 2D and 3D vectors and matrices"

std: (all units exist)

linalg
  linalg.vec2
    fn (x: f64, y: f64) -> vec2
    + constructs a 2D vector
    # construction
  linalg.vec3
    fn (x: f64, y: f64, z: f64) -> vec3
    + constructs a 3D vector
    # construction
  linalg.add3
    fn (a: vec3, b: vec3) -> vec3
    + returns component-wise sum
    # vector_arithmetic
  linalg.sub3
    fn (a: vec3, b: vec3) -> vec3
    + returns component-wise difference
    # vector_arithmetic
  linalg.scale3
    fn (v: vec3, k: f64) -> vec3
    + returns the vector scaled by k
    # vector_arithmetic
  linalg.dot3
    fn (a: vec3, b: vec3) -> f64
    + returns the scalar dot product
    # vector_arithmetic
  linalg.cross3
    fn (a: vec3, b: vec3) -> vec3
    + returns the right-handed cross product
    # vector_arithmetic
  linalg.length3
    fn (v: vec3) -> f64
    + returns the euclidean magnitude
    # vector_arithmetic
  linalg.normalize3
    fn (v: vec3) -> result[vec3, string]
    + returns v divided by its length
    - returns error when length is zero
    # vector_arithmetic
  linalg.mat3_identity
    fn () -> mat3
    + returns the 3x3 identity matrix
    # matrix_construction
  linalg.mat3_from_rows
    fn (r0: vec3, r1: vec3, r2: vec3) -> mat3
    + builds a 3x3 matrix from three row vectors
    # matrix_construction
  linalg.mat3_mul
    fn (a: mat3, b: mat3) -> mat3
    + returns standard 3x3 matrix product
    # matrix_arithmetic
  linalg.mat3_mul_vec3
    fn (m: mat3, v: vec3) -> vec3
    + transforms a vector by a 3x3 matrix
    # matrix_arithmetic
  linalg.mat3_transpose
    fn (m: mat3) -> mat3
    + returns the transpose
    # matrix_arithmetic
  linalg.mat3_determinant
    fn (m: mat3) -> f64
    + returns the determinant
    # matrix_arithmetic
  linalg.mat3_inverse
    fn (m: mat3) -> result[mat3, string]
    + returns the matrix inverse
    - returns error when the matrix is singular
    # matrix_arithmetic
