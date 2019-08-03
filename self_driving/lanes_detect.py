import cv2
import numpy as np
import matplotlib.pyplot as plt

def make_coords(image, line_params):
    slope, intercept = line_params
    y1 = image.shape[0]
    y2 = int(y1*3/5)
    x1 = int((y1 - intercept )/slope)
    x2 = int((y2 - intercept )/slope)
    return np.array([x1,y1,x2,y2])
  
def avg(image, lines):
    left = []
    right = []
    for line in lines:
        x1,y1, x2,y2 = line.reshape(4)
        params = np.polyfit((x1,x2), (y1,y2), 1)
        slope = params[0] 
        intercept = params[1]
        if slope <0:
            left.append((slope,intercept))
        else:
            right.append((slope,intercept))
    left_avg = np.average(left,axis=0)
    right_avg = np.average(right,axis=0)
    #print(left_avg,"left")
    #print(right_avg,"right")
    #[  -1.60826233 1200.57790627] left
    #[   1.00635503 -287.03395891] right
    left_line = make_coords(image, left_avg)
    right_line = make_coords(image, right_avg)
    return np.array([left_line,right_line])

def canny(image):
    gray_image = cv2.cvtColor(image, cv2.COLOR_RGB2GRAY)
    blur_image = cv2.GaussianBlur(gray_image, (5,5), 0)
    return cv2.Canny(blur_image, 15, 150)

def extract(image):
    height = image.shape[0]
    #width = image.shape[1]
    triangle = np.array([ (200,height), (1100, height), (550, 250) ])
    polygons = np.array([triangle])
    mask = np.zeros_like(image)
    region = cv2.fillPoly(mask, polygons, 255)
    masked_image = cv2.bitwise_and(image, region)
    return masked_image

def display_lines(image, lines):
    line_image = np.zeros_like(image)
    if lines is not None:
        for x1,y1, x2,y2 in lines:
            #x1,y1, x2,y2 = line.reshape(4)
            cv2.line(line_image,(x1,y1),(x2,y2),(255,0,0), 10)
    return line_image

#image = cv2.imread("test_image.jpg")
#lane_image = np.copy(image)
#canny_image = canny(lane_image)
#region = extract(canny_image)
#lines = cv2.HoughLinesP(region, 2, np.pi/180, 100, np.array([]), minLineLength=40, maxLineGap=5)
#avg_lines = avg(lane_image, lines)
#image_lines = display_lines(lane_image, avg_lines)
#combo_image = cv2.addWeighted(lane_image, 0.8, image_lines, 1, 1)
#############################################
#avg(combo_image,lines)
#cv2.imshow("result", combo_image)
#cv2.waitKey(0)
#plt.imshow(combo_image)
#plt.show()

cap = cv2.VideoCapture("test2.mp4")
while (cap.isOpened()):
    _, frame = cap.read()
    canny_image = canny(frame)
    region = extract(canny_image)
    lines = cv2.HoughLinesP(region, 2, np.pi/180, 100, np.array([]), minLineLength=40, maxLineGap=5)
    avg_lines = avg(frame, lines)
    image_lines = display_lines(frame, avg_lines)
    combo_image = cv2.addWeighted(frame, 0.8, image_lines, 1, 1)
    cv2.imshow("result", combo_image)
    if cv2.waitKey(1) & 0xFF == ord('q'):
        break
cap.release()
cv2.destroyAllWindows()
