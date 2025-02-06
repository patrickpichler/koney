// Copyright (c) 2025 Dynatrace LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package annotations

import (
	"encoding/json"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/dynatrace-oss/koney/api/v1alpha1"
	"github.com/dynatrace-oss/koney/internal/controller/constants"
	"github.com/dynatrace-oss/koney/internal/controller/utils"
)

const (
	testPodName   = "test-pod"
	testNamespace = "test-namespace"
	testCrdName   = "test-crd"
	testFilePath  = "/run/secrets/koney/service_token"
	testFileHash  = "75170fc230cd88f32e475ff4087f81d9"
)

var (
	deploymentStrategyValues = []string{
		"volumeMount",
		"containerExec",
		"kyvernoPolicy",
	}

	containersValues = [][]string{
		{},
		{"container1", "container2"},
	}

	trapTypeValues = []string{
		"filesystemHoneytoken",
	}

	changingFields = []string{"deploymentStrategy", "filePath", "fileContentHash", "readOnly"}

	annotationTraps []v1alpha1.Trap
)

// initializeTestTraps initializes the traps with all possible permutations of values to test the annotations
func initializeTestTraps() {
	for _, deploymentStrategy := range deploymentStrategyValues {
		for _, trapType := range trapTypeValues {
			switch trapType {
			case "filesystemHoneytoken":
				trap := v1alpha1.Trap{
					FilesystemHoneytoken: v1alpha1.FilesystemHoneytoken{
						FilePath:    testFilePath,
						FileContent: "someverysecrettoken",
						ReadOnly:    true,
					},
					DecoyDeployment: v1alpha1.DecoyDeployment{
						Strategy: deploymentStrategy,
					},
					MatchResources: v1alpha1.MatchResources{}, // This is not included in AnnotationTrap
				}
				annotationTraps = append(annotationTraps, trap)
			case "httpEndpoint":
				// TODO: Implement.
			case "httpPayload":
				// TODO: Implement.
			}
		}
	}
}

var _ = Describe("trapToAnnotationTrap", func() {
	Context("when converting a trap to an annotation trap", func() {
		It("should return an annotation trap with the same values", func() {
			for _, trap := range annotationTraps {
				for _, containers := range containersValues {
					annotationTrap, err := convertTrapToTrapAnnotation(trap, containers)
					Expect(err).ToNot(HaveOccurred())

					Expect(annotationTrap.DeploymentStrategy).To(Equal(trap.DecoyDeployment.Strategy))
					Expect(annotationTrap.TrapType()).To(Equal(trap.TrapType()))
					Expect(annotationTrap.Containers).To(Equal(containers))
					Expect(annotationTrap.CreatedAt).NotTo(BeEmpty())

					switch trap.TrapType() {
					case v1alpha1.FilesystemHoneytokenTrap:
						Expect(annotationTrap.FilesystemHoneytoken.FilePath).To(Equal(trap.FilesystemHoneytoken.FilePath))
					case v1alpha1.HttpEndpointTrap:
						// TODO: Implement.
					case v1alpha1.HttpPayloadTrap:
						// TODO: Implement.
					default:
						Fail("Unexpected trap type")
					}
				}
			}
		})
	})
})

var _ = Describe("AreTheSameTrap", func() {
	Context("when comparing a trap with an annotation trap created from it", func() {
		It("should return true", func() {
			for _, trap := range annotationTraps {
				// We manually craft the annotation trap
				var annotationTrap v1alpha1.TrapAnnotation
				switch trap.TrapType() {
				case v1alpha1.FilesystemHoneytokenTrap:
					annotationTrap = v1alpha1.TrapAnnotation{
						DeploymentStrategy: trap.DecoyDeployment.Strategy,
						Containers:         []string{}, // Not checked in the comparison
						CreatedAt:          "",         // Not checked in the comparison
						FilesystemHoneytoken: v1alpha1.FilesystemHoneytokenAnnotation{
							FilePath:        trap.FilesystemHoneytoken.FilePath,
							FileContentHash: utils.Hash(trap.FilesystemHoneytoken.FileContent),
							ReadOnly:        trap.FilesystemHoneytoken.ReadOnly,
						},
					}
				case v1alpha1.HttpEndpointTrap:
					// TODO: Implement.
				case v1alpha1.HttpPayloadTrap:
					// TODO: Implement.
				default:
					Fail("Unexpected trap type")
				}

				Expect(AreTheSameTrap(annotationTrap, trap)).To(BeTrue())
			}
		})
	})

	Context("when comparing a trap with an annotation trap with different values", func() {
		It("should return false", func() {
			// Then, we check the negative cases
			for _, trap := range annotationTraps {
				var annotationTrap v1alpha1.TrapAnnotation
				switch trap.TrapType() {
				case v1alpha1.FilesystemHoneytokenTrap:
					// We can modify the "deploymentStrategy" and the "filePath"
					for _, field := range changingFields {
						annotationTrap = v1alpha1.TrapAnnotation{
							DeploymentStrategy: trap.DecoyDeployment.Strategy,
							Containers:         []string{}, // Not checked in the comparison
							CreatedAt:          "",         // Not checked in the comparison
							FilesystemHoneytoken: v1alpha1.FilesystemHoneytokenAnnotation{
								FilePath:        trap.FilesystemHoneytoken.FilePath,
								FileContentHash: utils.Hash(trap.FilesystemHoneytoken.FileContent),
								ReadOnly:        trap.FilesystemHoneytoken.ReadOnly,
							},
						}
						// Modify the trap field
						switch field {
						case "deploymentStrategy":
							differentDeploymentStrategy := ""
							for _, strategy := range deploymentStrategyValues {
								if strategy != trap.DecoyDeployment.Strategy {
									differentDeploymentStrategy = strategy
									break
								}
							}
							annotationTrap.DeploymentStrategy = differentDeploymentStrategy
						case "filePath":
							annotationTrap.FilesystemHoneytoken.FilePath = fmt.Sprintf("%s/different", trap.FilesystemHoneytoken.FilePath)
						case "fileContentHash":
							annotationTrap.FilesystemHoneytoken.FileContentHash = testFileHash
						case "readOnly":
							annotationTrap.FilesystemHoneytoken.ReadOnly = !trap.FilesystemHoneytoken.ReadOnly
						}

						Expect(AreTheSameTrap(annotationTrap, trap)).To(BeFalse())
					}
				case v1alpha1.HttpEndpointTrap:
					// TODO: Implement.
				case v1alpha1.HttpPayloadTrap:
					// TODO: Implement.
				default:
					Fail("Unexpected trap type")
				}

				Expect(AreTheSameTrap(annotationTrap, trap)).To(BeFalse())
			}
		})
	})
})

var _ = Describe("trapToAnnotationTrap", func() {
	Context("when transforming a trap to an annotation trap", func() {
		It("should return an annotation trap with the same values", func() {
			for _, trap := range annotationTraps {
				for _, containers := range containersValues {
					annotationTrap, _ := convertTrapToTrapAnnotation(trap, containers)
					Expect(AreTheSameTrap(annotationTrap, trap)).To(BeTrue())
				}
			}
		})
	})
})

var _ = Describe("AddTrapToResourceAnnotations", func() {
	Context("when adding a trap to a pod annotations", func() {
		It("should add the trap to the annotations", func() {
			for _, trap := range annotationTraps {
				for _, containers := range containersValues {
					// We create a new test pod for each trap
					pod := corev1.Pod{
						ObjectMeta: metav1.ObjectMeta{
							Name:      testPodName,
							Namespace: testNamespace,
						},
					}
					// We add the trap to the pod annotations
					err := AddTrapToAnnotations(&pod, testCrdName, trap, containers)
					Expect(err).ToNot(HaveOccurred())

					// We check if the trap is in the annotations
					annotations := pod.Annotations[constants.AnnotationKeyChanges]

					// Unmarshal the annotations
					var annotationChanges []v1alpha1.ChangeAnnotation
					err = json.Unmarshal([]byte(annotations), &annotationChanges)
					Expect(err).ToNot(HaveOccurred())

					Expect(annotationChanges).To(HaveLen(1))
					annotation := annotationChanges[0]

					Expect(annotation.Traps).To(HaveLen(1))
					annotationTrap := annotation.Traps[0]

					Expect(annotation.DeceptionPolicyName).To(Equal(testCrdName))

					switch trap.TrapType() {
					case v1alpha1.FilesystemHoneytokenTrap:
						// Check that the trap in the annotations is the same as the trap we added
						Expect(annotationTrap.DeploymentStrategy).To(Equal(trap.DecoyDeployment.Strategy))
						Expect(annotationTrap.Containers).To(Equal(containers))
						Expect(annotationTrap.FilesystemHoneytoken.FilePath).To(Equal(trap.FilesystemHoneytoken.FilePath))
						Expect(annotationTrap.FilesystemHoneytoken.FileContentHash).To(Equal(utils.Hash(trap.FilesystemHoneytoken.FileContent)))
						Expect(annotationTrap.FilesystemHoneytoken.ReadOnly).To(Equal(trap.FilesystemHoneytoken.ReadOnly))
					case v1alpha1.HttpEndpointTrap:
						// TODO: Implement.
					case v1alpha1.HttpPayloadTrap:
						// TODO: Implement.
					default:
						Fail("Unexpected trap type")
					}
				}
			}
		})

		It("should update the trap in the annotations if it already exists", func() {
			for _, trap := range annotationTraps {
				// We create a single test pod for each trap
				// We reuse the same pod because we expect `AddTropToPodAnnotations`
				// to update the trap in the annotations if it already exists
				pod := corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      testPodName,
						Namespace: testNamespace,
					},
				}

				for _, containers := range containersValues {
					// We add the trap to the pod annotations
					err := AddTrapToAnnotations(&pod, testCrdName, trap, containers)
					Expect(err).ToNot(HaveOccurred())

					// We check if the trap is in the annotations
					annotations := pod.Annotations[constants.AnnotationKeyChanges]

					// Unmarshal the annotations
					var annotationChanges []v1alpha1.ChangeAnnotation
					err = json.Unmarshal([]byte(annotations), &annotationChanges)
					Expect(err).ToNot(HaveOccurred())

					Expect(annotationChanges).To(HaveLen(1))
					annotation := annotationChanges[0]

					Expect(annotation.Traps).To(HaveLen(1))
					annotationTrap := annotation.Traps[0]

					Expect(annotation.DeceptionPolicyName).To(Equal(testCrdName))

					switch trap.TrapType() {
					case v1alpha1.FilesystemHoneytokenTrap:
						// Check that the trap in the annotations is the same as the trap we added
						Expect(annotationTrap.DeploymentStrategy).To(Equal(trap.DecoyDeployment.Strategy))
						Expect(annotationTrap.Containers).To(Equal(containers))
						Expect(annotationTrap.FilesystemHoneytoken.FilePath).To(Equal(trap.FilesystemHoneytoken.FilePath))
					case v1alpha1.HttpEndpointTrap:
						// TODO: Implement.
					case v1alpha1.HttpPayloadTrap:
						// TODO: Implement.
					default:
						Fail("Unexpected trap type")
					}
				}
			}
		})

	})
})

var _ = Describe("UpdateContainersInTrapInResourceAnnotations", func() {
	Context("when updating the containers in a trap in the pod annotations", func() {
		It("should update the containers in the annotations", func() {
			for _, trap := range annotationTraps {
				for _, containers := range containersValues {
					// We create a new test pod for each trap
					pod := corev1.Pod{
						ObjectMeta: metav1.ObjectMeta{
							Name:      testPodName,
							Namespace: testNamespace,
						},
					}

					// We add the trap to the pod annotations
					err := AddTrapToAnnotations(&pod, testCrdName, trap, containers)
					Expect(err).ToNot(HaveOccurred())

					// We check if the trap is in the annotations
					annotations := pod.Annotations[constants.AnnotationKeyChanges]

					// Unmarshal the annotations
					var annotationChanges []v1alpha1.ChangeAnnotation
					err = json.Unmarshal([]byte(annotations), &annotationChanges)
					Expect(err).ToNot(HaveOccurred())

					// We perform some minor checks on the annotations
					Expect(annotationChanges).To(HaveLen(1))
					annotation := annotationChanges[0]

					Expect(annotation.Traps).To(HaveLen(1))
					annotationTrap := annotation.Traps[0]

					// We update the containers in the trap
					newContainers := append(containers, "some", "new", "containers")
					err = UpdateContainersInAnnotations(&pod, testCrdName, annotationTrap, newContainers)
					Expect(err).ToNot(HaveOccurred())

					// We check if the containers are updated in the annotations
					annotations = pod.Annotations[constants.AnnotationKeyChanges]

					// Unmarshal the annotations
					err = json.Unmarshal([]byte(annotations), &annotationChanges)
					Expect(err).ToNot(HaveOccurred())

					Expect(annotationChanges).To(HaveLen(1))

					annotation = annotationChanges[0]
					Expect(annotation.Traps).To(HaveLen(1))

					annotationTrap = annotation.Traps[0]
					Expect(annotationTrap.Containers).To(Equal(newContainers))
				}
			}
		})
	})
})

var _ = Describe("RemoveTrapFromResourceAnnotations", func() {
	Context("when removing a trap from the pod annotations", func() {
		It("should remove the trap from the annotations when there is a single trap", func() {
			for _, trap := range annotationTraps {
				// We create a new test pod for each trap
				pod := corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      testPodName,
						Namespace: testNamespace,
					},
				}

				// We add the trap to the pod annotations
				// We don't need to cycle the possible containers values for this test
				err := AddTrapToAnnotations(&pod, testCrdName, trap, containersValues[0])
				Expect(err).ToNot(HaveOccurred())

				// We check if the trap is in the annotations
				annotations := pod.Annotations[constants.AnnotationKeyChanges]

				// Unmarshal the annotations
				var annotationChanges []v1alpha1.ChangeAnnotation
				err = json.Unmarshal([]byte(annotations), &annotationChanges)
				Expect(err).ToNot(HaveOccurred())

				// We perform some minor checks on the annotations
				Expect(annotationChanges).To(HaveLen(1))
				annotation := annotationChanges[0]

				Expect(annotation.Traps).To(HaveLen(1))

				// We remove the trap from the annotations
				annotationTrap := annotation.Traps[0]
				err = RemoveTrapAnnotations(&pod, testCrdName, annotationTrap)
				Expect(err).ToNot(HaveOccurred())

				// We check if the trap is removed from the annotations
				annotations = pod.Annotations[constants.AnnotationKeyChanges]

				// We expect the annotations to be empty because we removed the only trap
				// and the AnnotationChanges annotation should be removed
				Expect(annotations).To(BeEmpty())
			}

		})
		It("should remove the trap from the annotations when there are multiple traps", func() {
			for _, trap1 := range annotationTraps {
				for _, trap2 := range annotationTraps {
					// Skip if the traps are the same
					// We check this based on the resulting annotation since
					// there is no Equals method for the Trap struct
					trap1Annotation, _ := convertTrapToTrapAnnotation(trap1, containersValues[0])
					trap2Annotation, _ := convertTrapToTrapAnnotation(trap2, containersValues[0])
					if trap1Annotation.Equals(&trap2Annotation, true) {
						continue
					}

					// We create a new test pod for each trap
					pod := corev1.Pod{
						ObjectMeta: metav1.ObjectMeta{
							Name:      testPodName,
							Namespace: testNamespace,
						},
					}

					// We add the traps to the pod annotations
					err := AddTrapToAnnotations(&pod, testCrdName, trap1, containersValues[0])
					Expect(err).ToNot(HaveOccurred())

					err = AddTrapToAnnotations(&pod, testCrdName, trap2, containersValues[0])
					Expect(err).ToNot(HaveOccurred())

					// We check if the traps are in the annotations
					annotations := pod.Annotations[constants.AnnotationKeyChanges]

					// Unmarshal the annotations
					var annotationChanges []v1alpha1.ChangeAnnotation
					err = json.Unmarshal([]byte(annotations), &annotationChanges)
					Expect(err).ToNot(HaveOccurred())

					// We perform some minor checks on the annotations
					Expect(annotationChanges).To(HaveLen(1))
					annotation := annotationChanges[0]

					// We remove the first trap from the annotations
					// Note: we use the AnnotationTrap extracted from the annotations
					//  because this is how the function is used in the controller
					//  (during the cleanup process, traps are extracted from the annotations)
					err = RemoveTrapAnnotations(&pod, testCrdName, annotation.Traps[0])
					Expect(err).ToNot(HaveOccurred())

					// We check if the trap is removed from the annotations
					annotations = pod.Annotations[constants.AnnotationKeyChanges]

					// Unmarshal the annotations
					err = json.Unmarshal([]byte(annotations), &annotationChanges)
					Expect(err).ToNot(HaveOccurred())

					Expect(annotationChanges).To(HaveLen(1))

					annotation = annotationChanges[0]
					Expect(annotation.Traps).To(HaveLen(1))
				}
			}
		})
	})
})

var _ = Describe("GetAnnotationChange", func() {
	Context("when getting the annotation changes for a specific DeceptionPolicy from a pod", func() {
		It("should return the annotation changes", func() {
			for _, trap := range annotationTraps {
				for _, containers := range containersValues {
					// We create a new test pod for each trap
					pod := corev1.Pod{
						ObjectMeta: metav1.ObjectMeta{
							Name:      testPodName,
							Namespace: testNamespace,
						},
					}

					// We add the trap to the pod annotations
					err := AddTrapToAnnotations(&pod, testCrdName, trap, containers)
					Expect(err).ToNot(HaveOccurred())

					// We get the annotation changes from the pod
					annotationChange, err := GetAnnotationChange(&pod, testCrdName)
					Expect(err).ToNot(HaveOccurred())

					Expect(annotationChange.DeceptionPolicyName).To(Equal(testCrdName))
					Expect(annotationChange.Traps).To(HaveLen(1))
					annotationTrap := annotationChange.Traps[0]

					switch trap.TrapType() {
					case v1alpha1.FilesystemHoneytokenTrap:
						Expect(annotationTrap.DeploymentStrategy).To(Equal(trap.DecoyDeployment.Strategy))
						Expect(annotationTrap.Containers).To(Equal(containers))
						Expect(annotationTrap.FilesystemHoneytoken.FilePath).To(Equal(trap.FilesystemHoneytoken.FilePath))
					case v1alpha1.HttpEndpointTrap:
						// TODO: Implement.
					case v1alpha1.HttpPayloadTrap:
						// TODO: Implement.
					default:
						Fail("Unexpected trap type")
					}
				}
			}
		})
	})
})
