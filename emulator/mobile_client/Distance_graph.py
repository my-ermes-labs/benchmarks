import matplotlib.pyplot as plt
import numpy as np
from geopy.distance import geodesic

# Function to calculate distances using the Haversine formula
def haversine_distance(coord1, coord2):
    return geodesic(coord1, coord2).km

# Define the points and nodes
latitudes = [28.998532, 34.415973, 41.212441, 46.498170, 47.813155]
nodes = [(31.728167, -96.717592), (37.195331, -99.592501), (42.065607, -94.485618), (44.546872, -101.516595)]
fixed_longitude = -98.652840
cloud_node = nodes[2]  # third node in the series

cloud_weigth = (latitudes[3]-latitudes[2]) / (latitudes[4]-latitudes[0])
edge_weigth = 1 - cloud_weigth

# Initialize distances
edge_cloud_distances = []
cloud_distances_for_edge = []
cloud_only_distances = []

# Compute distances for each km in the path
for lat in np.arange(latitudes[0], latitudes[-1] + 0.001, 0.009):
    # Distance to the cloud node (always calculated for Cloud-Only)
    distance_to_cloud = haversine_distance((lat, fixed_longitude), cloud_node)
    cloud_only_distances.append(distance_to_cloud)
    
    # Edge-Cloud specific distances
    if lat <= latitudes[1]:  # First segment
        distance_to_edge = haversine_distance((lat, fixed_longitude), nodes[0])
        edge_cloud_distances.append(distance_to_edge)
    elif lat <= latitudes[2]:  # Second segment
        distance_to_edge = haversine_distance((lat, fixed_longitude), nodes[1])
        edge_cloud_distances.append(distance_to_edge)
    elif lat <= latitudes[3]:  # Third segment (to Cloud node)
        cloud_distances_for_edge.append(distance_to_cloud)
    else:  # Fourth segment
        distance_to_edge = haversine_distance((lat, fixed_longitude), nodes[3])
        edge_cloud_distances.append(distance_to_edge)

# Calculate averages
average_edge_cloud_distance = np.mean(edge_cloud_distances)
average_cloud_distance_for_edge = np.mean(cloud_distances_for_edge)
average_cloud_only_distance = np.mean(cloud_only_distances)

# Data for the stacked bar chart with updated values and labels
edge_cloud = [average_edge_cloud_distance * edge_weigth, average_cloud_distance_for_edge * cloud_weigth]
cloud_only = [average_cloud_only_distance]

# Colors for the parts of the bars
colors_edge_cloud = ['#FFD8A8', '#B2F2BB']
colors_cloud_only = ['#B2F2BB']

# Create stacked bar chart
fig, ax = plt.subplots(figsize=(8, 6))

# Edge-Cloud bar
ax.bar('Edge-Cloud', edge_cloud[0], color=colors_edge_cloud[0], edgecolor='black', label='Edge')
ax.bar('Edge-Cloud', edge_cloud[1], bottom=edge_cloud[0], color=colors_edge_cloud[1], edgecolor='black', label='Cloud')

# Cloud-Only bar
ax.bar('Cloud-Only', cloud_only, color=colors_cloud_only, edgecolor='black')

# Labels and title
ax.set_xlabel('Setup')
ax.set_ylabel('Average Request Distance (KM)')
#ax.set_ylim(0, 600)

# Add legend
ax.legend()

# Display the bar graph
plt.show()

